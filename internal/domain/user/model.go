package user

import "time"

type User struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Role      string    `db:"role"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Banned    bool      `db:"banned"`
	CreatedAt time.Time `db:"created_at"`
}
