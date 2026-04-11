package profile

type ProfileResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateNameRequest struct {
	Name string `json:"name"`
}