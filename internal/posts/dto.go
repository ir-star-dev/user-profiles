package posts

import "user-profiles/internal/validator"

type CreateUserForm struct {
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	Name      string              `json:"name"`
	Role      string              `json:"role"`
	Validator validator.Validator `json:"-"`
}

type CreatePostForm struct {
	Title     string              `json:"title"`
	Content   string              `json:"content"`
	Excerpt   string              `json:"excerpt"`
	Validator validator.Validator `json:"-"`
}