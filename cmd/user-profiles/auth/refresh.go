package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

type refreshTokenService struct{}

func NewRefreshTokenService() RefreshTokenService {
	return &refreshTokenService{}
}

func (r *refreshTokenService) Generate() (string, error) {
	b := make([]byte, 32) // 256 bit
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (r *refreshTokenService) Hash(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	return hash[:]
}
