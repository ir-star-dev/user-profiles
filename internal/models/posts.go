package models

import (
	"html/template"
	"time"
)

type Post struct {
	Id        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	Excerpt   string    `db:"excerpt"`
	Slug      string    `db:"slug"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Approved  bool      `db:"approved"`
	UserId    int       `db:"user_id"`
}

type PostWithUserName struct {
	Id        int           `db:"id"`
	Title     string        `db:"title"`
	Content   template.HTML `db:"content"`
	Excerpt   string        `db:"excerpt"`
	Slug      string        `db:"slug"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
	Approved  bool          `db:"approved"`
	Username  string        `db:"username"`
	UserId    int           `db:"user_id"`
}

type Authors struct {
	Username string `db:"username"`
}

type PostsStatus struct {
	Published int `db:"published"`
	Pending   int `db:"pending"`
}

type UpdatePostRequest struct {
	Id        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	Excerpt   string    `db:"excerpt"`
	Approved  bool      `db:"approved"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

type PostRows struct {
	Id       int    `db:"id"`
	Title    string `db:"title"`
	Slug     string `db:"slug"`
	Approved bool   `db:"approved"`
	Username string `db:"username"`
}
