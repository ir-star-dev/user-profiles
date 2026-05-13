package auth

const (
	UserExists          = "User already exists"
	UserCreate          = "Cannot create a user"
	PassHash            = "Password error"
	LoginError          = "Wrong email"
	WrongPassword       = "Wrong password"
	NoRefreshtoken      = "No refresh token"
	InvalidRefreshToken = "Invalid refresh token"
	InvalidOrExpires	= "Invalid or expired refresh token"
	ReadingCookie		= "Error reading cookie"
	Unauthorized		= "Unauthorized"
	WrongCredentials    = "Wrong credentials"
	ReuseToken			= "Token reuse detected"
	SomethingWrong		= "Something went wrong. Try later, please!"
)