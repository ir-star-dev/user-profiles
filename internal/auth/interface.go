package auth

import (
	"database/sql"
	"user-profiles/internal/models"
	"user-profiles/internal/validator"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Register(email, password, name, role string) ([]validator.FormValidationErr, error)
	Login(email, password string) (*AuthResponse, []validator.FormValidationErr, error)
	Refresh(refreshToken string) (*AuthResponse, error)
	Logout(uId int, refreshToken string) error

	CreateLoginLog(uId int, agent, ip, device string) error
	GetLoginLog() ([]models.LoginsResponse, error)
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
	Save(token *models.RefreshToken) error
	Revoke(hash []byte, uId int) error
	RevokeFamily(familyId string) error
	FindTokenByHash(hash []byte) (*models.RefreshToken, error)
	FindUserIdByHash(hash []byte) (int, error)

	WithTx(fn func(repo RefreshRepository) error) error
}

type LoginsRepository interface {
	Save(login *models.Logins) error
	Get() (*[]models.LoginsResponse, error)
}

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	NamedExec(query string, arg any) (sql.Result, error)
	Get(dest any, query string, args ...any) error
}
