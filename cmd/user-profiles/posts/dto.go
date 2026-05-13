package posts

import "time"

type UpdatePostRequest struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Approved  bool      `json:"approved"`
	UpdatedAt time.Time `json:"updated_at"`
}

