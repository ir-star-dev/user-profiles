package auth_service

type Service interface {
	Register(email, password, name, role string) error
	// Login(email, password string) (*AuthResponse, error)
	// Refresh(refreshToken string) (*AuthResponse, error)
	// Logout(uId int, refreshToken string) error
}