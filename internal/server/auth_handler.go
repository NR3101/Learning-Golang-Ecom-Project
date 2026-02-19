package server

import (
	"github.com/NR3101/go-ecom-project/internal/dto"
	"github.com/NR3101/go-ecom-project/internal/utils"
	"github.com/gin-gonic/gin"
)

// register handles user registration requests
func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	response, err := s.authService.Register(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to register user", err)
		return
	}

	utils.CreatedResponse(c, "User registered successfully", response)
}

// login handles user login requests
func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	response, err := s.authService.Login(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid email or password")
		return
	}

	utils.SuccessResponse(c, "User logged in successfully", response)
}

// refreshToken handles token refresh requests
func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	response, err := s.authService.RefreshToken(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid refresh token")
		return
	}

	utils.SuccessResponse(c, "Token refreshed successfully", response)
}

// logout handles user logout requests
func (s *Server) logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	err := s.authService.Logout(req.RefreshToken)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to logout user", err)
		return
	}

	utils.SuccessResponse(c, "User logged out successfully", nil)
}
