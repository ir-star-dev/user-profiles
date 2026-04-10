package token

type Repository interface {
	Save(token *RefreshToken) error
	FindTokenByHash(hash []byte) (*RefreshToken, error)
	FindUserIdByHash(hash []byte) (int, error)
	Revoke(hash []byte, uId int) error
}

