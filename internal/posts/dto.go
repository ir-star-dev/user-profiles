package posts

type UpdatePostRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Approved bool   `json:"approved"`
}
