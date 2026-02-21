package server

import (
	"github.com/NR3101/go-ecom-project/internal/dto"
	"github.com/NR3101/go-ecom-project/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param registerRequest body dto.RegisterRequest true "User registration data"
// @Success 201 {object} utils.Response{data=dto.AuthResponse} "User registered successfully"
// @Failure 400 {object} utils.Response "Invalid request body or user already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/register [post]
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

// @Summary Login user
// @Description Login with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginRequest body dto.LoginRequest true "User sign-in data"
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "User logged in successfully"
// @Failure 400 {object} utils.Response "Invalid request body"
// @Failure 401 {object} utils.Response "Invalid email or password"
// @Router /auth/login [post]
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

// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param refreshTokenRequest body dto.RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "Token refreshed successfully"
// @Failure 400 {object} utils.Response "Invalid request body"
// @Failure 401 {object} utils.Response "Invalid refresh token"
// @Router /auth/refresh [post]
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

// @Summary Logout user
// @Description Logout user and invalidate refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param logoutRequest body dto.RefreshTokenRequest true "Refresh token to invalidate"
// @Success 200 {object} utils.Response "User logged out successfully"
// @Failure 400 {object} utils.Response "Invalid request body"
// @Failure 500 {object} utils.Response "Failed to logout user"
// @Router /auth/logout [post]
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
