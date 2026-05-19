package posts

import "errors"

var (
	PostNotFound   = errors.New("Posts not found")
	LongTitle      = "Title length more than 99 characters"
	EmptyTitle     = "Title cannot be blank"
	ShortExcerpt   = "Excerpt length less than 100 characters"
	LongExcerpt    = "Excerpt length more than 300 characters"
	EmptyExcerpt   = "Excerpt cannot be blank"
	EmptyContent   = "Content cannot be blank"
	InvalidForm    = "Invalid form"
	SomethingWrong = "Something went wrong. Try later, please!"
)
