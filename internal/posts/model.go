package posts

import "time"

type Post struct {
	Id        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	Approved  bool      `db:"approved"`
	UserId    int       `db:"user_id"`
}
