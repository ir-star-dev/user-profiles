package posts

import "time"

type Post struct {
	Id        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Approved  bool      `db:"approved"`
	UserId    int       `db:"user_id"`
}

type PostWithUserName struct {
	Id        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Approved  bool      `db:"approved"`
	UserName  string    `db:"user_name"`
}