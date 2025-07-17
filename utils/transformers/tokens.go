package transformers

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
)

func GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateTokenFromString(input string) string {
	hasher := crypto.SHA256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GenerateMD5Hash(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}
