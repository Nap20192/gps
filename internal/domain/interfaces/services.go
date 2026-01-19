package interfaces

import "gps/internal/app_services/auth"

type AuthService interface {
	SignUp(username, password string) error
	Login(username, password string) (string, error)
	ValidateToken(tokenString string) (*auth.JWTClaims, error)
}
