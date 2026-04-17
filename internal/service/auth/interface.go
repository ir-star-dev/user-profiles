package auth_service

import auth_dto "user-profiles/internal/dto/auth"

type Service interface {
	Register(email, password, name, role string) error
	Login(email, password string) (*auth_dto.AuthResponse, error)
	// Refresh(refreshToken string) (*AuthResponse, error)
	// Logout(uId int, refreshToken string) error
}
