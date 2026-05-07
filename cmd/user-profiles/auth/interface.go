package auth

import (
	"database/sql"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Register(email, password, name, role string) error
	Login(email, password string) (*AuthResponse, error)
	Refresh(refreshToken string) (*AuthResponse, error)
	Logout(uId int, refreshToken string) error
}

type JWTService interface {
	Create(uId int, role string, ban *bool) (string, error)
	Parse(token string) (jwt.MapClaims, error)
}

type RefreshTokenService interface {
	Generate() (string, error)
	Hash(token string) []byte
}

type RefreshRepository interface {
	Save(token *RefreshToken) error
	Revoke(hash []byte, uId int) error
	RevokeFamily(familyId string) error
	FindTokenByHash(hash []byte) (*RefreshToken, error)
	FindUserIdByHash(hash []byte) (int, error)

	WithTx(fn func(repo RefreshRepository) error) error
}

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	NamedExec(query string, arg any) (sql.Result, error)
	Get(dest any, query string, args ...any) error
}