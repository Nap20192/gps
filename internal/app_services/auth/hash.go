package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
)

const (
	passwordSaltSize   = 16
	passwordIterations = 1000
)

func hashPassword(password string) (salt string, hash string, err error) {
	saltBytes := make([]byte, passwordSaltSize)

	if _, err = rand.Read(saltBytes); err != nil {
		return "", "", err
	}

	hashBytes := sha256.Sum256(append(saltBytes, []byte(password)...))

	for i := 1; i < passwordIterations; i++ {
		hashBytes = sha256.Sum256(hashBytes[:])
	}

	salt = base64.RawStdEncoding.EncodeToString(saltBytes)
	hash = base64.RawStdEncoding.EncodeToString(hashBytes[:])
	return
}

func verifyPassword(password, salt, expectedHash string) bool {
	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return false
	}

	expectedBytes, err := base64.RawStdEncoding.DecodeString(expectedHash)
	if err != nil {
		return false
	}

	hashBytes := sha256.Sum256(append(saltBytes, []byte(password)...))

	for i := 1; i < passwordIterations; i++ {
		hashBytes = sha256.Sum256(hashBytes[:])
	}

	return subtle.ConstantTimeCompare(hashBytes[:], expectedBytes) == 1
}
