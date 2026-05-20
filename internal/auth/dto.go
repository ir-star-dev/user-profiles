package auth

type AuthResponse struct {
	UserId            int
	Access            string
	Refresh           string
	RefreshHash       []byte
}
