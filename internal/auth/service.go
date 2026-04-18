package auth

import (
	"errors"
	"time"
	"user-profiles/internal/users"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type authService struct {
	uRepo     users.Repository
	tRepo     RefreshRepository
	jService  JWTService
	rtService RefreshTokenService
}

func NewAuthService(uRepo users.Repository, tRepo RefreshRepository, jS JWTService, rtS RefreshTokenService) Service {
	return &authService{
		uRepo:     uRepo,
		tRepo:     tRepo,
		jService:  jS,
		rtService: rtS,
	}
}

func (s *authService) Register(email, password, name, role string) error {
	
	existedUser, _ := s.uRepo.FindByEmail(email)
	if existedUser != nil {
		return errors.New(UserExists)
	}
	password = strings.TrimSpace(password)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := &users.UserWithRole{
		Email:    strings.TrimSpace(email),
		Name:     strings.TrimSpace(name),
		Password: string(hash),
		Role:     role,
	}
	_, err = s.uRepo.Create(newUser)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) Login(email, password string) (*AuthResponse, error) {
	existedUser, _ := s.uRepo.FindByEmail(email)
	if existedUser == nil {
		return nil, errors.New(LoginError)
	}
	password = strings.TrimSpace(password)
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return nil, errors.New(WrongPassword)
	}
	tokens, err := s.generateTokens(existedUser.Id, existedUser.Role)
	if err != nil {
		return nil, err
	}
	fId := uuid.NewString()
	err = s.saveRefreshToken(tokens.RefreshHash, existedUser.Id, fId)
	if err != nil {
		return nil, err
	}
	res := &AuthResponse{
		Access: tokens.Access,
		Refresh: tokens.Refresh,
		RefreshHash: tokens.RefreshHash,
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

// func (s *authService) Refresh(refreshToken string) (*AuthResponse, error) {
// 	t, err := s.validateRefresh(refreshToken)
// 	if err != nil {
// 		return nil, errors.New(InvalidOrExpires)
// 	}
// 	tokens, err := s.generateTokens(t.UserId, "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	newRefreshHash := s.rtService.Hash(Refresh)

// 	// 🔥 transaction
// 	err = s.tRepo.WithTx(func(repo token.Repository) error {
// 		if !t.Revoked {
// 			// 💥 compromise detected
// 			_ = repo.RevokeFamily(t.FamilyID)
// 			if err := repo.Revoke(t.TokenHash, t.UserId); err != nil {
// 				return err
// 			}
// 		}

// 		return repo.Save(&token.RefreshToken{
// 			TokenHash: newRefreshHash,
// 			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
// 			Revoked:   false,
// 			UserId:    t.UserId,
// 			FamilyID:  t.FamilyID,
// 		})
// 	})

// 	if err != nil {
// 		return nil, err
// 	}
// 	return tokens, nil
// }

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

// func (s *authService) validateRefresh(token string) (*token.RefreshToken, error) {
// 	hash := s.rtService.Hash(token)
// 	existedToken, err := s.tRepo.FindTokenByHash(hash)
// 	if existedToken == nil {
// 		return nil, err
// 	}
// 	if existedToken.Revoked || existedToken.ExpiresAt.Before(time.Now()) {
// 		return nil, err
// 	}
// 	return existedToken, nil
// }
