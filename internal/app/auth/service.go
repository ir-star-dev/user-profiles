package auth

import (
	"time"
	"errors"
	"user-profiles/internal/domain/token"
	"user-profiles/internal/domain/user"
	"user-profiles/internal/infrastructure/security"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(email, password, name, host string) (*RegisterResponse, error)
	Login(email, password, host string) (*AuthResponse, error)
	Refresh(refreshToken, host string) (*AuthResponse, error)
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

func (s *authService) Register(email, password, name, host string) (*RegisterResponse, error) {
	existedUser, _ := s.uRepo.FindByEmail(email)
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
	_, err = s.uRepo.Create(newUser)
	if err != nil {
		return nil, err
	}
	tokens, err := s.generateTokens(host, existedUser.Id)
	if err != nil {
		return nil, err
	}
	res := &RegisterResponse{
		Token: tokens.Access,
	}
	return res, nil
}

func (s *authService) Login(email, password, host string) (*AuthResponse, error) {
	existedUser, _ := s.uRepo.FindByEmail(email)
	if existedUser == nil {
		return nil, errors.New(LoginError)
	}
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return nil, errors.New(WrongPassword)
	}
	tokens, err := s.generateTokens(host, existedUser.Id)

	err = s.saveRefreshToken(tokens.RefreshHash, existedUser.Id)
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

func (s *authService) Refresh(refreshToken, host string) (*AuthResponse, error) {
	t, err := s.validateRefresh(refreshToken)
	if err != nil { 
        return nil, errors.New(InvalidOrExpires)
    }
	tokens, err := s.generateTokens(host, t.UserId)
	if err != nil {
		return nil, err
	}
	newRefreshHash := s.rtService.Hash(tokens.Refresh)
	err = s.saveRefreshToken(newRefreshHash, t.UserId)
	if err != nil {
		return nil, err
	}	
	err = s.rotate(t.TokenHash, t.UserId)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}


func (s *authService) saveRefreshToken(hash []byte, uId int) error {
	token := &token.RefreshToken{
		TokenHash: hash,
		ExpiresAt: time.Now().Add(7*24*time.Hour),
		Revoked: false,
		UserId: uId,
	}

	err := s.tRepo.Save(token)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) generateTokens(host string, uId int) (*AuthResponse, error) {
	// Access token lived 15 min
	jwt, err := s.jService.Create(host, uId)
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

func (s *authService) rotate(oldHash []byte, uId int) error {
	err := s.revoke(oldHash, uId)
	if err != nil {
		return err
	}
	return nil
}