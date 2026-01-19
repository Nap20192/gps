package auth

import (
	"context"
	"fmt"

	"gps/internal/config"
	"gps/internal/domain/models"

	"github.com/google/uuid"
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

type Input struct {
	Username string `json:"username" binding:"required,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type repo interface {
	CreateUser(ctx context.Context, username, passwordHash string) (uuid.UUID, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
}

type AuthService struct {
	repo repo
}

func NewAuthService(repo repo) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) SignUp(ctx context.Context, input Input) (string, error) {
	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return "", err
	}

	id, err := s.repo.CreateUser(ctx, input.Username, hashedPassword)
	if err != nil {
		return "", err
	}
	token, err := generateToken(JWTClaims{
		UserID:   id,
		Username: input.Username,
	})

	return token, nil
}

func (s *AuthService) LogIn(ctx context.Context, input Input) (string, error) {

	inputUsername := input.Username
	user, err := s.repo.GetUserByUsername(ctx, inputUsername)

	if err != nil {
		return "", ErrUserNotFound
	}

	isValid := verifyPassword(input.Password, user.PasswordHash)

	if !isValid {
		return "", ErrInvalidCredentials
	}

	token, err := generateToken(JWTClaims{
		UserID:   user.UserID,
		Username: user.Username,
	})

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) ParseToken(tokenStr string) (JWTClaims, error) {
	return parseToken(tokenStr)
}
