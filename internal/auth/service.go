package auth

import (
	"errors"
	"time"
	"user-profiles/internal/models"
	"user-profiles/internal/users"
	"user-profiles/internal/validator"

	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	uRepo     users.Repository
	tRepo     RefreshRepository
	lRepo     LoginsRepository
	jService  JWTService
	rtService RefreshTokenService
}

func NewAuthService(uRepo users.Repository, tRepo RefreshRepository, lRepo LoginsRepository, jS JWTService, rtS RefreshTokenService) AuthService {
	return &authService{
		uRepo:     uRepo,
		tRepo:     tRepo,
		lRepo:     lRepo,
		jService:  jS,
		rtService: rtS,
	}
}

func (s *authService) Register(email, password, name, role string) ([]validator.FormValidationErr, error) {
	var formValiErr []validator.FormValidationErr
	var v validator.Validator
	v.CheckField(validator.Matches(email, validator.EmailRX), "email", EmailNotValid)
	v.CheckField(validator.MaxChars(password, 8), "password", LongPass)
	v.CheckField(validator.NotBlank(password), "password", EmptyPass)
	v.CheckField(validator.NotBlank(name), "name", EmptyName)
	v.CheckField(validator.PermittedValue(role, "user"), "role", WrongRole)

	if !v.Valid() {
		for key, value := range v.FieldErrors {
			formValiErr = append(formValiErr, validator.FormValidationErr{
				Name:    key,
				Message: value,
			})
		}
		return formValiErr, errors.New(InvalidForm)
	}

	existedUser, err := s.uRepo.FindByEmail(email)
	if existedUser != nil {
		formValiErr = append(formValiErr, validator.FormValidationErr{
			Name:    "exist",
			Message: UserExists,
		})
		return formValiErr, errors.New(UserExists)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		formValiErr = append(formValiErr, validator.FormValidationErr{
			Name:    "password",
			Message: PassHash,
		})
		return formValiErr, errors.New(PassHash)
	}

	username, _, _ := strings.Cut(email, "@")

	newUser := models.UserWithRole{
		Email:    email,
		Name:     name,
		Username: username,
		Password: string(hash),
		Role:     role,
	}
	_, err = s.uRepo.Create(&newUser)
	if err != nil {
		formValiErr = append(formValiErr, validator.FormValidationErr{
			Name:    "form",
			Message: SomethingWrong,
		})
		return formValiErr, errors.New(SomethingWrong)
	}
	return nil, nil
}

func (s *authService) Login(email, password string) (*AuthResponse, []validator.FormValidationErr, error) {
	var formValiErr []validator.FormValidationErr
	existedUser, err := s.uRepo.FindByEmail(email)
	if existedUser == nil {
		formValiErr = append(formValiErr, validator.FormValidationErr{
			Name:    "email",
			Message: LoginError,
		})
		return nil, formValiErr, errors.New(LoginError)
	}
	password = strings.TrimSpace(password)
	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		formValiErr = append(formValiErr, validator.FormValidationErr{
			Name:    "password",
			Message: WrongPassword,
		})
		return nil, formValiErr, errors.New(WrongPassword)
	}
	tokens, err := s.generateTokens(existedUser.Id, existedUser.Role)
	if err != nil {
		formValiErr = append(formValiErr, validator.FormValidationErr{
			Name:    "form",
			Message: SomethingWrong,
		})
		return nil, formValiErr, errors.New(SomethingWrong)
	}
	fId := uuid.NewString()
	err = s.saveRefreshToken(tokens.RefreshHash, existedUser.Id, fId)
	if err != nil {
		formValiErr = append(formValiErr, validator.FormValidationErr{
			Name:    "form",
			Message: SomethingWrong,
		})
		return nil, formValiErr, errors.New(SomethingWrong)
	}
	res := &AuthResponse{
		UserId:            existedUser.Id,
		Access:            tokens.Access,
		Refresh:           tokens.Refresh,
		RefreshHash:       tokens.RefreshHash,
	}
	return res, nil, nil
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

			return repo.Save(&models.RefreshToken{
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
	token := &models.RefreshToken{
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

func (s *authService) validateRefresh(token string) (*models.RefreshToken, error) {
	hash := s.rtService.Hash(token)
	existedToken, err := s.tRepo.FindTokenByHash(hash)
	if existedToken == nil {
		return nil, err
	}
	return existedToken, nil
}

func (s *authService) CreateLoginLog(uId int, agent, ip, device string) error {
	logins := models.Logins{
		UserId:    uId,
		UserAgent: agent,
		Ip:        ip,
		Device:    device,
	}
	err := s.lRepo.Save(&logins)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) GetLoginLog() ([]models.LoginsResponse, error) {
	logs, err := s.lRepo.Get()
	if err != nil {
		return nil, err
	}
	return *logs, nil
}