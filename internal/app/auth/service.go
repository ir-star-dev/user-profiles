package auth

import (
	"errors"
	"time"
	"user-profiles/internal/domain/token"
	"user-profiles/internal/domain/user"
	"user-profiles/internal/infrastructure/security"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(email, password, name string) (*RegisterResponse, error)
	Login(email, password string) (*AuthResponse, error)
	Refresh(refreshToken string) (*AuthResponse, error)
	Logout(uId int, refreshToken string) error
}

type authService struct {
	uRepo user.Repository
	tRepo token.Repository
	jService security.JWTService
	rtService security.RefreshTokenService
}

func NewAuthService(uRepo user.Repository, tRepo token.Repository, jS security.JWTService, rtS security.RefreshTokenService) Service {
	return &authService {
		uRepo: uRepo,
		tRepo: tRepo,
		jService: jS,
		rtService: rtS,
	}
}

func (s *authService) Register(email, password, name string) (*RegisterResponse, error) {
	existedUser, err := s.uRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existedUser != nil {
		return nil, errors.New(UserExists)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUser := &user.User{
		Email:    email,
		Name:     name,
		Password: hash,
	}
	user, err := s.uRepo.Create(newUser)
	if err != nil {
		return nil, err
	}
	tokens, err := s.generateTokens(user.Id)
	if err != nil {
		return nil, err
	}
	res := &RegisterResponse{
		Token: tokens.Access,
	}
	return res, nil
}

func (s *authService) Login(email, password string) (*AuthResponse, error) {
	existedUser, err := s.uRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existedUser == nil {
		return nil, errors.New(LoginError)
	}
	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return nil, errors.New(WrongPassword)
	}
	tokens, err := s.generateTokens(existedUser.Id)
	if err != nil {
		return nil, err
	}
	fId := uuid.NewString()
	err = s.saveRefreshToken(tokens.RefreshHash, existedUser.Id, fId)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *authService) Logout(uId int, refreshToken string) error {
	tokenHash := s.rtService.Hash(refreshToken)
	err := s.revoke(tokenHash, uId)
	if err != nil {
		return errors.New(Unauthorized)
	}
	return nil
}

func (s *authService) Refresh(refreshToken string) (*AuthResponse, error) {
	t, err := s.validateRefresh(refreshToken)
	if err != nil { 
        return nil, errors.New(InvalidOrExpires)
    }
	if t.Revoked {
		// 💥 compromise detected
		_ = s.tRepo.RevokeFamily(t.FamilyID)
		return nil, errors.New(ReuseToken)
	}

	tokens, err := s.generateTokens(t.UserId)
	if err != nil {
		return nil, err
	}
	newRefreshHash := s.rtService.Hash(tokens.Refresh)
	err = s.saveRefreshToken(newRefreshHash, t.UserId, t.FamilyID)
	if err != nil {
		return nil, err
	}	
	err = s.revoke(t.TokenHash, t.UserId)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}


func (s *authService) saveRefreshToken(hash []byte, uId int, fId string) error {
	token := &token.RefreshToken{
		TokenHash: hash,
		ExpiresAt: time.Now().Add(7*24*time.Hour),
		Revoked: false,
		UserId: uId,
		FamilyID:  fId,
	}

	err := s.tRepo.Save(token)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) generateTokens(uId int) (*AuthResponse, error) {
	// Access token lived 15 min
	jwt, err := s.jService.Create(uId)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.rtService.Generate()
	if err != nil {
		return nil, err
	}
	hash := s.rtService.Hash(refreshToken)

	return &AuthResponse{
		Access: jwt,
		Refresh: refreshToken,
		RefreshHash: hash,
	}, nil
}

func (s *authService) revoke(hash []byte, uId int) error {
	err := s.tRepo.Revoke(hash, uId)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) validateRefresh(token string) (*token.RefreshToken, error) {
	hash := s.rtService.Hash(token)
	existedToken, err := s.tRepo.FindTokenByHash(hash)
	if existedToken == nil { 
        return nil, err
    }
	if existedToken.Revoked || existedToken.ExpiresAt.Before(time.Now()) {
		return nil, err
	}
	return existedToken, nil
}
