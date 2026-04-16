package auth_service

const (
	UserExists          = "User already exists"
	UserCreate          = "Cannot create a user"
	PassHash            = "Password error"
	LoginError          = "Wrong email or password"
	WrongPassword       = "Wrong password"
	NoRefreshtoken      = "No refresh token"
	InvalidRefreshToken = "Invalid refresh token"
	InvalidOrExpires	= "Invalid or expired refresh token"
	ReadingCookie		= "Error reading cookie"
	Unauthorized		= "Unauthorized"
	WrongCredentials    = "Wrong credentials"
	ReuseToken			= "Token reuse detected"
)