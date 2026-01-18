package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const (
	userIDContextKey contextKey = "user_id"
)

var ErrUnauthorized = errors.New("unauthorized")

type AuthClient interface {
	ValidateToken(ctx context.Context, token string) (string, error)
}

type BaseAuthClient struct {
	baseURL      string
	validatePath string
	disabled     bool
	client       *http.Client
}

func NewBaseAuthClient(baseURL, validatePath string, disabled bool) *BaseAuthClient {
	if validatePath == "" {
		validatePath = "/auth/validate"
	}
	return &BaseAuthClient{
		baseURL:      strings.TrimRight(baseURL, "/"),
		validatePath: validatePath,
		disabled:     disabled,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *BaseAuthClient) ValidateToken(ctx context.Context, token string) (string, error) {
	if c.disabled {
		return "", nil
	}
	if c.baseURL == "" {
		return "", ErrUnauthorized
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+c.validatePath, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", ErrUnauthorized
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("auth service error")
	}

	var payload struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	return payload.UserID, nil
}

func AuthMiddleware(authClient AuthClient) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userID, err := authClient.ValidateToken(ctx, token)
			if err != nil {
				status := http.StatusInternalServerError
				if errors.Is(err, ErrUnauthorized) {
					status = http.StatusUnauthorized
				}
				http.Error(w, "Unauthorized", status)
				return
			}

			ctx = context.WithValue(ctx, userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
