package users

type UpdateNameRequest struct {
	Name *string `json:"name"`
}

type StatUsers struct {
	Count int
	Role  string
}
