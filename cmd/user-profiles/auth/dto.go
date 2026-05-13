package auth

import (
	"user-profiles/cmd/user-profiles/panel"
	"user-profiles/internal/validator"
)

type RegisterInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Validator validator.Validator `json:"-"`
}

type AuthResponse struct {
	UserId            int
	Access            string
	Refresh           string
	RefreshHash       []byte
	FormValidationErr []panel.FormValidationErr
}
