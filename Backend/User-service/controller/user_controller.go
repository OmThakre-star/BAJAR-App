package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"user-service/models"
	"user-service/utils"

	"go.uber.org/zap"
)

// GetProfile godoc
// @Summary Get user profile
// @Tags User
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/me [get]
func (uc *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value("userID")
	var uid string
	switch v := val.(type) {
	case string:
		uid = v
	case int:
		uid = strconv.Itoa(v)
	default:
		uc.Logger.Error("Invalid user id in context", zap.Any("userID", val))
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	user, err := uc.UserService.GetUserByID(r.Context(), uid)
	if err != nil {
		uc.Logger.Error("User not found", zap.String("userID", uid), zap.Error(err))
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	user.Password = ""
	uc.Logger.Info("User profile fetched", zap.String("userID", uid))
	json.NewEncoder(w).Encode(user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Tags User
// @Accept json
// @Produce json
// @Param profile body models.User true "Profile update info"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/me [put]
func (uc *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		uc.Logger.Error("UpdateProfile: userID missing in context")
		http.Error(w, "userID missing", http.StatusBadRequest)
		return
	}
	var req struct {
		Name   string `json:"name"`
		Phone  string `json:"phone"`
		DOB    string `json:"dob"`
		Gender string `json:"gender"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.Logger.Error("UpdateProfile: invalid JSON", zap.Error(err))
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	user, err := uc.UserService.GetUserByID(r.Context(), userID)
	if err != nil {
		uc.Logger.Error("UpdateProfile: user not found", zap.String("userID", userID), zap.Error(err))
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	user.Name = req.Name
	user.Phone = req.Phone
	// Optionally parse and set DOB, Gender
	if err := uc.UserService.UpdateUser(r.Context(), user); err != nil {
		uc.Logger.Error("UpdateProfile: update failed", zap.Error(err))
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	uc.Logger.Info("User profile updated", zap.String("userID", userID))
	json.NewEncoder(w).Encode(user)
}

// ChangePassword godoc
// @Summary Change user password
// @Tags User
// @Accept json
// @Produce json
// @Param password body models.ChangePasswordRequest true "Old and new password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/change-password [post]
func (uc *UserController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		uc.Logger.Error("ChangePassword: userID missing in context")
		http.Error(w, "userID missing", http.StatusBadRequest)
		return
	}

	var ChangePasswordRequest models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&ChangePasswordRequest); err != nil {
		uc.Logger.Error("ChangePassword: invalid JSON", zap.Error(err))
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	user, err := uc.UserService.GetUserByID(r.Context(), userID)
	if err != nil {
		uc.Logger.Error("ChangePassword: user not found", zap.String("userID", userID), zap.Error(err))
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	if err := utils.CheckPassword(user.Password, ChangePasswordRequest.OldPassword); err != nil {
		uc.Logger.Error("ChangePassword: old password incorrect", zap.Error(err))
		http.Error(w, "old password incorrect", http.StatusUnauthorized)
		return
	}
	hash, err := utils.HashPassword(ChangePasswordRequest.NewPassword)
	if err != nil {
		uc.Logger.Error("ChangePassword: hash failed", zap.Error(err))
		http.Error(w, "hash failed", http.StatusInternalServerError)
		return
	}
	user.Password = hash
	if err := uc.UserService.UpdateUser(r.Context(), user); err != nil {
		uc.Logger.Error("ChangePassword: update failed", zap.Error(err))
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	uc.Logger.Info("Password changed", zap.String("userID", userID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed"})
}

// DeleteAccount godoc
// @Summary Delete user account
// @Tags User
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/delete-account [delete]
func (uc *UserController) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		uc.Logger.Error("DeleteAccount: userID missing in context")
		http.Error(w, "userID missing", http.StatusBadRequest)
		return
	}
	if err := uc.UserService.DeleteUser(r.Context(), userID); err != nil {
		uc.Logger.Error("DeleteAccount: delete failed", zap.Error(err))
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	uc.Logger.Info("Account deleted", zap.String("userID", userID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Account deleted"})
}

// ListUsers godoc
// @Summary List all users
// @Tags Admin
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /admin/users [get]
func (uc *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uc.UserService.ListUsers(r.Context())
	if err != nil {
		uc.Logger.Error("Could not fetch users", zap.Error(err))
		http.Error(w, "could not fetch users", http.StatusInternalServerError)
		return
	}
	for i := range users {
		users[i].Password = ""
	}
	uc.Logger.Info("User list fetched", zap.Int("count", len(users)))
	json.NewEncoder(w).Encode(users)
}

// AdminListUsers godoc
// @Summary Admin: List users with filters
// @Tags Admin
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param role query string false "User role"
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /admin/users [get]
func (uc UserController) AdminListUsers(w http.ResponseWriter, r *http.Request) {
	// Example: parse pagination/filter params from query
	// page := r.URL.Query().Get("page")
	// limit := r.URL.Query().Get("limit")
	// role := r.URL.Query().Get("role")
	// TODO: Convert page/limit to int, apply filters in service/repo
	users, err := uc.UserService.ListUsers(r.Context())
	if err != nil {
		uc.Logger.Error("AdminListUsers: fetch failed", zap.Error(err))
		http.Error(w, "could not fetch users", http.StatusInternalServerError)
		return
	}
	// TODO: Apply pagination/filter logic here
	json.NewEncoder(w).Encode(users)
}

// AdminGetUser godoc
// @Summary Admin: Get user by ID
// @Tags Admin
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /admin/users/{id} [get]
func (uc UserController) AdminGetUser(w http.ResponseWriter, r *http.Request) {
	// Example: get user ID from URL
	id := r.URL.Path[len("/admin/users/"):] // crude extraction, use mux vars in real code
	user, err := uc.UserService.GetUserByID(r.Context(), id)
	if err != nil {
		uc.Logger.Error("AdminGetUser: user not found", zap.String("userID", id), zap.Error(err))
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// AdminUpdateRole godoc
// @Summary Admin: Update user role
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param role body models.UpdateRoleRequest true "New role"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /admin/users/{id}/role [put]
func (uc UserController) AdminUpdateRole(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/admin/users/") : len(r.URL.Path)-len("/role")]
	var req models.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.Logger.Error("AdminUpdateRole: invalid JSON", zap.Error(err))
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	user, err := uc.UserService.GetUserByID(r.Context(), id)
	if err != nil {
		uc.Logger.Error("AdminUpdateRole: user not found", zap.String("userID", id), zap.Error(err))
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	user.Role = req.Role
	if err := uc.UserService.UpdateUser(r.Context(), user); err != nil {
		uc.Logger.Error("AdminUpdateRole: update failed", zap.Error(err))
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	uc.Logger.Info("User role updated", zap.String("userID", id), zap.String("role", req.Role))
	json.NewEncoder(w).Encode(user)
}

// AdminDeleteUser godoc
// @Summary Admin: Delete/deactivate user
// @Tags Admin
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /admin/users/{id} [delete]
func (uc UserController) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/admin/users/"):] // crude extraction, use mux vars in real code
	if err := uc.UserService.DeleteUser(r.Context(), id); err != nil {
		uc.Logger.Error("AdminDeleteUser: delete failed", zap.Error(err))
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	uc.Logger.Info("User deactivated/banned", zap.String("userID", id))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deactivated/banned"})
}
