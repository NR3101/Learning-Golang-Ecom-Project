package resolver

import (
	"context"
	"errors"

	"github.com/NR3101/go-ecom-project/internal/utils"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

const (
	adminRole = "admin"
)

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	userID := ctx.Value(utils.UserID)
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
	userRole := ctx.Value(utils.UserRole)
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

// getPageAndLimit extracts page and limit from pointers, applying defaults and bounds
func getPageAndLimit(page, limit *int) (int, int) {
	p := 1
	l := 10
	if page != nil {
		p = *page
	}
	if limit != nil {
		l = *limit
	}
	if p <= 0 {
		p = 1
	}
	if l <= 0 {
		l = 10
	}
	return p, l
}
