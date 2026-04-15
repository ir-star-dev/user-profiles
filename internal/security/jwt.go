package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	Create(uId int, role string, ban *bool) (string, error)
	Parse(token string) (jwt.MapClaims, error)
}

type JWT struct {
	Secret string
}

func NewJWTService(secret string) JWTService {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(uId int, role string, ban *bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":    jwt.NewNumericDate(time.Now()),
		"exp":    jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		"sub":    uId,
		"role":   role,
		"banned": ban,
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
