package service

import (
	"context"
	"fmt"
	"time"
	"user-service/models"
	"user-service/repository"
	"user-service/utils"

	"go.uber.org/zap"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.repo.CreateUser(ctx, user)
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.repo.UpdateUser(ctx, user)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(ctx, id)
}

// ListUsers returns all users
func (s *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.ListUsers(ctx)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

// AuthenticateUser checks user credentials
func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
	return s.repo.GetUserByEmailAndPassword(ctx, email, password)
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, id string, newPassword string) error {
	return s.repo.ChangePassword(ctx, id, newPassword)
}

// ResetPassword resets a user's password
func (s *UserService) ResetPassword(ctx context.Context, email, newPassword string) error {
	return s.repo.ResetPassword(ctx, email, newPassword)
}

// Register registers a new user
func (s *UserService) Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error) {
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hash,
		Phone:    req.Phone,
		Address:  req.Address,
	}
	if err := s.repo.CreateUser(ctx, &user); err != nil {
		return nil, err
	}
	// TODO: send verification email/OTP here if needed
	resp := &models.RegisterResponse{ID: user.ID, Email: user.Email}
	return resp, nil
}

// Login logs in a user
func (s *UserService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return nil, err
	}
	accessToken, _ := utils.CreateToken(user.ID, 15*time.Minute)
	refreshToken, _ := utils.CreateToken(user.ID, 7*24*time.Hour)
	resp := &models.LoginResponse{ID: user.ID, Email: user.Email, Token: accessToken}
	// Optionally, add refreshToken to resp struct or return as a separate value
	// For now, just log it (in production, send as httpOnly cookie or in response)
	zap.L().Info("Refresh token generated", zap.String("refreshToken", refreshToken))
	return resp, nil
}

// Logout logs out a user
func (s *UserService) Logout(ctx context.Context, userID string) error {
	zap.L().Info("Logout called", zap.String("userID", userID))
	// TODO: Implement token invalidation logic
	return nil
}

// Refresh refreshes the JWT token
func (s *UserService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.ParseToken(refreshToken)
	if err != nil {
		return "", err
	}
	// Optionally, check if token is blacklisted or revoked
	userID := claims.Subject
	// Generate new access token
	// You may want to check user existence or status here
	// Convert userID to uint if needed
	// For simplicity, assume userID is valid
	var uid uint
	_, err = fmt.Sscanf(userID, "%d", &uid)
	if err != nil {
		return "", err
	}
	newAccessToken, err := utils.CreateToken(uid, 15*time.Minute)
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}

// ForgotPassword handles the password reset request
func (s *UserService) ForgotPassword(ctx context.Context, email string) error {
	zap.L().Info("ForgotPassword called", zap.String("email", email))
	// TODO: Implement password reset request logic (send OTP/email)
	return nil
}

// VerifyEmail verifies the user's email
func (s *UserService) VerifyEmail(ctx context.Context, tokenOrOTP string) error {
	zap.L().Info("VerifyEmail called", zap.String("tokenOrOTP", tokenOrOTP))
	// TODO: Implement email verification logic
	return nil
}
