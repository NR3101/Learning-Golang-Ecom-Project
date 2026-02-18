package server

import (
	"github.com/NR3101/go-ecom-project/internal/dto"
	"github.com/NR3101/go-ecom-project/internal/services"
	"github.com/NR3101/go-ecom-project/internal/utils"
	"github.com/gin-gonic/gin"
)

// getProfile retrieves the profile information of the authenticated user.
func (s *Server) getProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	userService := services.NewUserService(s.db)

	profile, err := userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "User profile retrieved successfully", profile)
}

// updateProfile updates the profile information of the authenticated user.
func (s *Server) updateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	userService := services.NewUserService(s.db)
	updatedProfile, err := userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update user profile", err)
		return
	}

	utils.SuccessResponse(c, "User profile updated successfully", updatedProfile)
}
