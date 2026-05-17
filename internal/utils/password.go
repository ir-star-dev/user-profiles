package utils

import (
	"crypto/rand"
	"math/big"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GeneratePassword(length int) (string, error) {
	password := make([]byte, length)

	for i := range password {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}

		password[i] = chars[n.Int64()]
	}

	return string(password), nil
}
