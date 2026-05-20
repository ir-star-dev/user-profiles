package auth

const (
	UserExists     = "Email already taken"
	PassHash       = "Password error"
	LoginError     = "Wrong email"
	WrongPassword  = "Wrong password"
	Unauthorized   = "Unauthorized"
	SomethingWrong = "Something went wrong. Try later, please!"
	EmailNotValid  = "Email is not valid"
	LongPass       = "Password length more than 8 characters"
	EmptyPass      = "Password cannot be blank"
	EmptyName      = "Name cannot be blank"
	WrongRole      = "Role can be only a user"
	InvalidForm    = "Invalid form"
	InvalidToken   = "Invalid token"
)
