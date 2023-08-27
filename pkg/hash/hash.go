package hash

import (
	"crypto/sha256"
	"fmt"
)

type PasswordHash interface {
	Hash(password string) string
}

type SHA256 struct {
	salt string
}

func NewPasswordHashSHA256(salt string) *SHA256 {
	return &SHA256{salt: salt}
}

func (sha *SHA256) Hash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(sha.salt)))
}
