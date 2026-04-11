package token

type Repository interface {
	Save(token *RefreshToken) error
	Revoke(hash []byte, uId int) error
	RevokeFamily(familyId string) error
	FindTokenByHash(hash []byte) (*RefreshToken, error)
	FindUserIdByHash(hash []byte) (int, error)
	
	WithTx(fn func(repo Repository) error) error
}

