package profile

import "time"

type ProfileResponse struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Banned bool   `json:"banned"`
}

type ProfileResponseForAdmin struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	Banned    bool      `json:"banned"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateNameRequest struct {
	Name *string `json:"name"`
}
