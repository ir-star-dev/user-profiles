package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	Create(host string, uId int) (string, error)
	Parse(token string) (jwt.MapClaims, error)
}

type JWT struct {
	Secret string
}

type Claims struct {
	Audience string
	Issuer   string
}

func NewJWTService(secret string) JWTService {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(host string, uId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		"iss": host,
		"sub": uId,
		//"jti":   "uniq_token_id",
		"aud": host,
	})

	s, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return s, nil
}

func (j *JWT) Parse(tokenStr string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("Invalid token")
	}
	return claims, nil
}
