package resolver

import (
	"fmt"
	"strconv"

	"github.com/NR3101/go-ecom-project/internal/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Dependency injection for services
type Resolver struct {
	authService    services.AuthServiceInterface
	userService    services.UserServiceInterface
	productService services.ProductServiceInterface
	cartService    services.CartServiceInterface
	orderService   services.OrderServiceInterface
}

// NewResolver creates a new Resolver with the provided services.
func NewResolver(
	authService services.AuthServiceInterface,
	userService services.UserServiceInterface,
	productService services.ProductServiceInterface,
	cartService services.CartServiceInterface,
	orderService services.OrderServiceInterface,
) *Resolver {
	return &Resolver{
		authService:    authService,
		userService:    userService,
		productService: productService,
		cartService:    cartService,
		orderService:   orderService,
	}
}

func (r *Resolver) parseID(id string) (uint, error) {
	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %w", err)
	}
	return uint(parsedID), nil
}
