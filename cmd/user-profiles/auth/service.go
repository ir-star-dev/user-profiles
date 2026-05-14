package auth

import (
	"errors"
	"time"
	"user-profiles/cmd/user-profiles/panel"
	"user-profiles/cmd/user-profiles/users"
	"user-profiles/internal/validator"

	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	uRepo     users.Repository
	tRepo     RefreshRepository
	jService  JWTService
	rtService RefreshTokenService
}

func NewAuthService(uRepo users.Repository, tRepo RefreshRepository, jS JWTService, rtS RefreshTokenService) AuthService {
	return &authService{
		uRepo:     uRepo,
		tRepo:     tRepo,
		jService:  jS,
		rtService: rtS,
	}
}

func (s *authService) Register(email, password, name, role string) (*AuthResponse, error) {
	var formValiErr []panel.FormValidationErr
	form := RegisterInput{
		Email:     email,
		Name:      name,
		Role:      role,
		Password:  password,
		Validator: validator.Validator{},
	}

	form.Validator.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Email is not valid")
	form.Validator.CheckField(validator.MaxChars(form.Password, 8), "password", "Password length more than 8 characters")
	form.Validator.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")
	form.Validator.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be blank")
	form.Validator.CheckField(validator.PermittedValue(form.Role, "user"), "role", "Role can be user")

	if !form.Validator.Valid() {
		for key, value := range form.Validator.FieldErrors {
			formValiErr = append(formValiErr, panel.FormValidationErr{
				Name:    key,
				Message: value,
			})
		}
		return &AuthResponse{
			FormValidationErr: formValiErr,
		}, errors.New("Invalid form")
	}

	existedUser, err := s.uRepo.FindByEmail(email)
	if existedUser != nil {
		formValiErr = append(formValiErr, panel.FormValidationErr{
			Name:    "exist",
			Message: UserExists,
		})
		return &AuthResponse{
			FormValidationErr: formValiErr,
		}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		formValiErr = append(formValiErr, panel.FormValidationErr{
			Name:    "password",
			Message: PassHash,
		})
		return &AuthResponse{
			FormValidationErr: formValiErr,
		}, err
	}

	username, _, _ := strings.Cut(email, "@")

	newUser := users.UserWithRole{
		Email:    email,
		Name:     name,
		Username: username,
		Password: string(hash),
		Role:     role,
	}
	_, err = s.uRepo.Create(&newUser)
	if err != nil {
		formValiErr = append(formValiErr, panel.FormValidationErr{
			Name:    "form",
			Message: SomethingWrong,
		})
		return &AuthResponse{
			FormValidationErr: formValiErr,
		}, err
	}
	return &AuthResponse{}, nil
}

func (s *authService) Login(email, password string) (*AuthResponse, error) {
	var formValiErr []panel.FormValidationErr
	existedUser, err := s.uRepo.FindByEmail(email)
	if existedUser == nil {
		formValiErr = append(formValiErr, panel.FormValidationErr{
			Name:    "email",
			Message: LoginError,
		})
		return &AuthResponse{
			FormValidationErr: formValiErr,
		}, err
	}
	password = strings.TrimSpace(password)
	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		formValiErr = append(formValiErr, panel.FormValidationErr{
			Name:    "password",
			Message: WrongPassword,
		})
		return &AuthResponse{
			FormValidationErr: formValiErr,
		}, err
	}
	tokens, err := s.generateTokens(existedUser.Id, existedUser.Role)
	if err != nil {
		formValiErr = append(formValiErr, panel.FormValidationErr{
			Name:    "form",
			Message: SomethingWrong,
		})
		return &AuthResponse{
			FormValidationErr: formValiErr,
		}, err
	}
	fId := uuid.NewString()
	err = s.saveRefreshToken(tokens.RefreshHash, existedUser.Id, fId)
	if err != nil {
		formValiErr = append(formValiErr, panel.FormValidationErr{
			Name:    "form",
			Message: SomethingWrong,
		})
		return &AuthResponse{
			FormValidationErr: formValiErr,
		}, err
	}
	res := &AuthResponse{
		UserId:            existedUser.Id,
		Access:            tokens.Access,
		Refresh:           tokens.Refresh,
		RefreshHash:       tokens.RefreshHash,
		FormValidationErr: []panel.FormValidationErr{},
	}
	return res, nil
}

func (s *authService) Logout(uId int, refreshToken string) error {
	tokenHash := s.rtService.Hash(refreshToken)
	err := s.tRepo.Revoke(tokenHash, uId)
	if err != nil {
		return errors.New(Unauthorized)
	}
	return nil
}

func (s *authService) Refresh(refreshToken string) (*AuthResponse, error) {
	t, err := s.validateRefresh(refreshToken)
	if err != nil {
		return nil, err
	}
	if !t.Revoked || !t.ExpiresAt.Before(time.Now()) {
		tokens, err := s.generateTokens(t.UserId, "")
		if err != nil {
			return nil, err
		}
		newRefreshHash := s.rtService.Hash(tokens.Refresh)
		// 🔥 transaction
		err = s.tRepo.WithTx(func(repo RefreshRepository) error {
			if !t.Revoked {
				// 💥 compromise detected
				_ = repo.RevokeFamily(t.FamilyID)
				if err := repo.Revoke(t.TokenHash, t.UserId); err != nil {
					return err
				}
			}

			return repo.Save(&RefreshToken{
				TokenHash: newRefreshHash,
				ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
				Revoked:   false,
				UserId:    t.UserId,
				FamilyID:  t.FamilyID,
			})
		})
		if err != nil {
			return nil, err
		}
		return tokens, nil
	}
	return nil, err
}

func (s *authService) saveRefreshToken(hash []byte, uId int, fId string) error {
	token := &RefreshToken{
		TokenHash: hash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
		UserId:    uId,
		FamilyID:  fId,
	}

	err := s.tRepo.Save(token)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) generateTokens(uId int, role string) (*AuthResponse, error) {
	// Access token lived 15 min
	if role == "" {
		dbRole, err := s.uRepo.FindRoleByUserId(uId)
		if err != nil {
			return nil, err
		}
		role = dbRole
	}
	ban, err := s.uRepo.FindBanStatus(uId)
	if err != nil {
		return nil, err
	}
	jwt, err := s.jService.Create(uId, role, ban)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.rtService.Generate()
	if err != nil {
		return nil, err
	}
	hash := s.rtService.Hash(refreshToken)

	return &AuthResponse{
		Access:      jwt,
		Refresh:     refreshToken,
		RefreshHash: hash,
	}, nil
}

func (s *authService) validateRefresh(token string) (*RefreshToken, error) {
	hash := s.rtService.Hash(token)
	existedToken, err := s.tRepo.FindTokenByHash(hash)
	if existedToken == nil {
		return nil, err
	}
	return existedToken, nil
}
