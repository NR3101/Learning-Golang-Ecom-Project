package services

import (
	"errors"
	"log/slog"
	"time"

	"github.com/NR3101/go-ecom-project/internal/config"
	"github.com/NR3101/go-ecom-project/internal/dto"
	"github.com/NR3101/go-ecom-project/internal/events"
	"github.com/NR3101/go-ecom-project/internal/models"
	"github.com/NR3101/go-ecom-project/internal/utils"
	"gorm.io/gorm"
)

const (
	ErrEmailAlreadyExists = "email already exists"
)

type AuthService struct {
	db             *gorm.DB
	config         *config.Config
	eventPublisher events.Publisher
}

func NewAuthService(db *gorm.DB, config *config.Config, eventPublisher events.Publisher) *AuthService {
	return &AuthService{
		db:             db,
		config:         config,
		eventPublisher: eventPublisher,
	}
}

// Register creates a new user account and returns an authentication response with tokens.
func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if email already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New(ErrEmailAlreadyExists)
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create the user
	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.UserRoleCustomer,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Create user cart
	cart := models.Cart{
		UserID: user.ID,
	}
	if err := s.db.Create(&cart).Error; err != nil {
		slog.Error("Failed to create cart for user", "user_id", user.ID, "error", err)
	}

	// Generate tokens and return response
	return s.generateAuthResponse(&user)
}

// Login authenticates a user and returns an authentication response with tokens.
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	var user models.User
	if err := s.db.Where("email = ? AND is_Active = ?", req.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens and return response
	return s.generateAuthResponse(&user)
}

// RefreshToken validates the provided refresh token, and if valid, generates new access and refresh tokens.
func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(req.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Find if refresh token exists in DB and is not expired
	var refreshToken models.RefreshToken
	if err := s.db.Where("token = ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	// Find user
	var user models.User
	if err := s.db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Delete old refresh token
	s.db.Delete(&refreshToken)

	// Generate new tokens and return response
	return s.generateAuthResponse(&user)
}

// Logout invalidates the provided refresh token by deleting it from the database.
func (s *AuthService) Logout(refreshToken string) error {
	// Delete refresh token from DB
	result := s.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("refresh token not found or already logged out")
	}
	return nil
}

// generateAuthResponse is a helper function that generates access and refresh tokens for a user.
func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	// Generate tokens
	accessToken, refreshToken, err := utils.GenerateToken(
		&s.config.JWT,
		user.ID,
		user.Email,
		string(user.Role))
	if err != nil {
		return nil, err
	}

	// Save refresh token in DB
	refreshTokenRecord := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpiresIn),
	}
	if err := s.db.Create(&refreshTokenRecord).Error; err != nil {
		return nil, err
	}

	// Prepare response
	userResponse := dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	err = s.eventPublisher.Publish("user_authenticated", userResponse, map[string]string{})
	if err != nil {
		slog.Error("Failed to publish user_authenticated event", "user_id", user.ID, "error", err)
		return nil, err
	}

	return &dto.AuthResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
