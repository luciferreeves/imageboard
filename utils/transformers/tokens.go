package transformers

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateTokenFromString(input string) string {
	hasher := crypto.SHA256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}
