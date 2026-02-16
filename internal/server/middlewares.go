package server

import (
	"strings"

	"github.com/NR3101/go-ecom-project/internal/models"
	"github.com/NR3101/go-ecom-project/internal/utils"
	"github.com/gin-gonic/gin"
)

// authMiddleware is a Gin middleware function that checks for the presence of a valid JWT token.
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header and validate it is not empty
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Validate the token format (should be "Bearer <token>")
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid Authorization header format")
			c.Abort()
			return
		}

		// Validate the token and extract claims
		tokenString := tokenParts[1]
		claims, err := utils.ValidateToken(tokenString, s.config.JWT.Secret)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Store user information in the context for use in handlers
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		// Proceed to the next handler
		c.Next()
	}
}

// adminMiddleware is a Gin middleware function that checks if the authenticated user has an admin role.
func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			utils.UnauthorizedResponse(c, "Admin access required")
			c.Abort()
			return
		}

		if role != string(models.UserRoleAdmin) {
			utils.UnauthorizedResponse(c, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
