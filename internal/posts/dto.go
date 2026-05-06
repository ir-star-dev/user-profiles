package posts

import "time"

type UpdatePostRequest struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Approved  bool      `json:"approved"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
