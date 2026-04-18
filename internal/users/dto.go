package users

import "time"

type UsersProfileResponse struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Banned bool   `json:"banned"`
}

type FullProfileResponseForAdmin struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	Banned    bool      `json:"banned"`
	CreatedAt time.Time `json:"created_at"`
}

type ShortProfileResponseForAdmin struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
}

type UpdateNameRequest struct {
	Name *string `json:"name"`
}
