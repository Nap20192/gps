package middleware

import (
	"context"
	"errors"
)

var ErrNoUserInContext = errors.New("unauthorized: invalid user context")

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	if !ok {
		return "", ErrNoUserInContext
	}
	return userID, nil
}
