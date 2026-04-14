package user

import "time"

type User struct {
	Id        int       `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Name      string    `db:"name"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}
