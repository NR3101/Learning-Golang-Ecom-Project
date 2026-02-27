package resolver

import (
	"context"
	"errors"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

const (
	adminRole = "admin"
)

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return 0, ErrUnauthorized
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, ErrUnauthorized
	}

	return id, nil
}

func GetUserRoleFromContext(ctx context.Context) (string, error) {
	userRole := ctx.Value("user_role")
	if userRole == nil {
		return "", ErrUnauthorized
	}

	role, ok := userRole.(string)
	if !ok {
		return "", ErrUnauthorized
	}

	return role, nil
}

func IsAdmin(ctx context.Context) bool {
	role, err := GetUserRoleFromContext(ctx)
	if err != nil {
		return false
	}

	return role == adminRole
}
