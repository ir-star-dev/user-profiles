package models

import "time"

type User struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Username  string    `db:"username"`
	RoleId    int       `db:"role_id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Banned    bool      `db:"banned"`
	CreatedAt time.Time `db:"created_at"`
}

type UserWithRole struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Username  string    `db:"username"`
	Role      string    `db:"role"`
	Banned    bool      `db:"banned"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

type StatUsers struct {
	Count  int    `db:"count"`
	Role   string `db:"role"`
	Banned int    `db:"banned"`
}

type Roles struct {
	Id   int    `db:"id"`
	Role string `db:"role"`
}
