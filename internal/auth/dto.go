package auth

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=12"`
	Name     string `json:"name" validate:"required,min=3"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
}

type AuthResponse struct {
	UserId		int
	Access      string
	Refresh     string
	RefreshHash []byte
}

type AuthViewError struct {
	Message string
}