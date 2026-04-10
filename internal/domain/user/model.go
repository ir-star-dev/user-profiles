package user

type User struct {
	Id       int    `db:"id"`
	Email    string `db:"email"`
	Password []byte `db:"password"`
	Name     string `db:"name"`
}
