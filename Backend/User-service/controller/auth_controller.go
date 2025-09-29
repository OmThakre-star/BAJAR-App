package controller

import (
	"encoding/json"
	"net/http"
	"user-service/models"
	"user-service/service"

	"go.uber.org/zap"
)

type UserController struct {
	UserService service.UserService
	Logger      *zap.Logger
}

// Register godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration info"
// @Success 200 {object} models.RegisterResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.Logger.Error("Invalid JSON in Register", zap.Error(err))
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	resp, err := uc.UserService.Register(r.Context(), &req)
	if err != nil {
		uc.Logger.Error("User registration failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uc.Logger.Info("User registered", zap.Uint("userID", resp.ID), zap.String("email", resp.Email))
	json.NewEncoder(w).Encode(resp)
}

// Login godoc
// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.LoginRequest true "User login info"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.Logger.Error("Invalid JSON in Login", zap.Error(err))
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	resp, err := uc.UserService.Login(r.Context(), &req)
	if err != nil {
		uc.Logger.Error("User login failed", zap.Error(err), zap.String("email", req.Email))
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	uc.Logger.Info("User logged in", zap.Uint("userID", resp.ID), zap.String("email", resp.Email))
	json.NewEncoder(w).Encode(resp)
}

// Logout godoc
// @Summary Logout user
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /auth/logout [post]
func (uc UserController) Logout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID == nil {
		uc.Logger.Error("Logout: userID missing in context")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "userID missing"})
		return
	}
	if err := uc.UserService.Logout(r.Context(), userID.(string)); err != nil {
		uc.Logger.Error("Logout failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	uc.Logger.Info("User logged out", zap.String("userID", userID.(string)))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}

// Refresh godoc
// @Summary Refresh JWT token
// @Tags Auth
// @Produce json
// @Param Authorization header string true "Refresh token"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (uc UserController) Refresh(w http.ResponseWriter, r *http.Request) {
	// Example: get refresh token from request (body/header/cookie)
	refreshToken := r.Header.Get("Authorization")
	token, err := uc.UserService.Refresh(r.Context(), refreshToken)
	if err != nil {
		uc.Logger.Error("Token refresh failed", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	uc.Logger.Info("Token refreshed")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// ForgotPassword godoc
// @Summary Request password reset
// @Tags Auth
// @Accept json
// @Produce json
// @Param email body string true "User email"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/forgot-password [post]
func (uc UserController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.Logger.Error("ForgotPassword: invalid JSON", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}
	if err := uc.UserService.ForgotPassword(r.Context(), req.Email); err != nil {
		uc.Logger.Error("ForgotPassword failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	uc.Logger.Info("Password reset requested", zap.String("email", req.Email))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset requested"})
}

// ResetPassword godoc
// @Summary Reset password
// @Tags Auth
// @Accept json
// @Produce json
// @Param reset body string true "Reset token/OTP and new password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/reset-password [post]
func (uc UserController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TokenOrOTP  string `json:"token_or_otp"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.Logger.Error("ResetPassword: invalid JSON", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}
	if err := uc.UserService.ResetPassword(r.Context(), req.TokenOrOTP, req.NewPassword); err != nil {
		uc.Logger.Error("ResetPassword failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	uc.Logger.Info("Password reset", zap.String("token_or_otp", req.TokenOrOTP))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset"})
}

// VerifyEmail godoc
// @Summary Verify user email
// @Tags Auth
// @Accept json
// @Produce json
// @Param verify body string true "Verification token/OTP"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/verify-email [post]
func (uc UserController) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TokenOrOTP string `json:"token_or_otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.Logger.Error("VerifyEmail: invalid JSON", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}
	if err := uc.UserService.VerifyEmail(r.Context(), req.TokenOrOTP); err != nil {
		uc.Logger.Error("VerifyEmail failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	uc.Logger.Info("Email verified", zap.String("token_or_otp", req.TokenOrOTP))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email verified"})
}
