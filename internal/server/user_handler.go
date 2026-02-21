package server

import (
	"github.com/NR3101/go-ecom-project/internal/dto"
	"github.com/NR3101/go-ecom-project/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get user profile
// @Description Get the profile information of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.UserResponse} "User profile retrieved successfully"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users/profile [get]
func (s *Server) getProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	profile, err := s.userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "User profile retrieved successfully", profile)
}

// @Summary Update user profile
// @Description Update the profile information of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param updateProfileRequest body dto.UpdateProfileRequest true "User profile update data"
// @Success 200 {object} utils.Response{data=dto.UserResponse} "User profile updated successfully"
// @Failure 400 {object} utils.Response "Invalid request body"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users/profile [put]
func (s *Server) updateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	updatedProfile, err := s.userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update user profile", err)
		return
	}

	utils.SuccessResponse(c, "User profile updated successfully", updatedProfile)
}
