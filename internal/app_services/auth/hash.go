package auth

import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) ( hash string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func verifyPassword(password, expectedHash string) bool {
	if expectedHash == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(expectedHash), []byte(password)) == nil
}
