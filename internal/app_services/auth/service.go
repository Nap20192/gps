package auth

import (
	"fmt"

	"gps/internal/config"
)

var (
	ErrUserNotFound       = fmt.Errorf("user not found")
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
	ErrFailedToCreateUser = fmt.Errorf("failed to create user")
	ErrUnauthorized       = fmt.Errorf("unauthorized")
)

func ConfigureJWT(cfg config.JWTConfig) {
	SetJWTConfig(cfg.Secret, cfg.Expiry)
}

// type RegisterInput struct {
// 	Attrs    map[string]string `json:"attrs,omitempty"`
// 	Email    string            `json:"email" binding:"required,email,max=64"`
// 	Password string            `json:"password" binding:"required,min=8,max=64"`
// }

// type LoginInput struct {
// 	Email    string `json:"email" binding:"required,email,max=64"`
// 	Password string `json:"password" binding:"required,min=8,max=64"`
// }

// type AuthService struct {
// 	queries sqlc.Queries
// }

// func NewAuthService(queries sqlc.Queries) *AuthService {
// 	return &AuthService{
// 		queries: queries,
// 	}
// }

// func (s *AuthService) SignUp(ctx context.Context, input RegisterInput, role core.UserRole) (string, error) {
// 	salt, hashedPassword, err := hashPassword(input.Password)
// 	if err != nil {
// 		return "", err
// 	}

// 	user, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
// 		Email:        input.Email,
// 		Role:         role.String(),
// 		Status:       core.UserStatusInactive.String(),
// 		PasswordHash: hashedPassword,
// 		Salt:         salt,
// 		Attrs:        input.Attrs,
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	token, err := generateToken(JWTClaims{
// 		UserID: user.ID,
// 		Email:  user.Email,
// 		Role:   user.Role,
// 	})

// 	return token, nil
// }

// func (s *AuthService) LogIn(ctx context.Context, input LoginInput) (string, error) {
// 	inputEmail := input.Email

// 	user, err := s.queries.GetUserByEmail(ctx, inputEmail)
// 	if err != nil {
// 		return "", ErrUserNotFound
// 	}

// 	isValid := verifyPassword(input.Password, user.Salt, user.PasswordHash)

// 	if !isValid {
// 		return "", ErrInvalidCredentials
// 	}

// 	token, err := generateToken(JWTClaims{
// 		UserID: user.ID,
// 		Email:  user.Email,
// 		Role:   user.Role,
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	return token, nil
// }

// func (s *AuthService) ParseToken(tokenStr string) (JWTClaims, error) {
// 	return parseToken(tokenStr)
// }
