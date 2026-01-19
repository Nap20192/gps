package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID   uuid.UUID
	Username string
	jwt.RegisteredClaims
}

var (
	tokenTTL  = time.Hour * 1
	jwtSecret = []byte("super-secret-key")
)

func InitJwt(secret string, ttl time.Duration) {
	tokenTTL = ttl
	jwtSecret = []byte(secret)
}

func SetJWTConfig(secret string, ttl time.Duration) {
	if secret != "" {
		jwtSecret = []byte(secret)
	}
	if ttl > 0 {
		tokenTTL = ttl
	}
}

func generateToken(claims JWTClaims) (string, error) {
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func parseToken(tokenString string) (JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		},
	)
	if err != nil {
		return JWTClaims{}, err
	}

	if !token.Valid {
		return JWTClaims{}, errors.New("invalid token")
	}

	return *claims, nil
}
