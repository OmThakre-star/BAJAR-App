package router

import (
	"net/http"
	"user-service/controller"
	"user-service/middleware"

	"github.com/gorilla/mux"
)

type Route struct {
	Path      string
	Method    string
	Handler   func(http.ResponseWriter, *http.Request)
	Protected bool
}

func NewUserRouter(userController *controller.UserController) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.RateLimitMiddleware)

	routes := []Route{
		{"/auth/register", "POST", userController.Register, false},
		{"/auth/login", "POST", userController.Login, false},
		{"/auth/logout", "POST", userController.Logout, true},
		{"/auth/refresh", "POST", userController.Refresh, false},
		{"/auth/forgot-password", "POST", userController.ForgotPassword, false},
		{"/auth/reset-password", "POST", userController.ResetPassword, false},
		{"/auth/verify-email", "POST", userController.VerifyEmail, false},
		{"/api/me", "GET", userController.GetProfile, true},
		{"/api/me", "PUT", userController.UpdateProfile, true},
		{"api/delete-account", "DELETE", userController.DeleteAccount, true},
		{"/api/change-password", "POST", userController.ChangePassword, true},
		{"/admin/users", "GET", userController.AdminListUsers, true},
		{"/admin/users/{id}", "GET", userController.AdminGetUser, true},
		{"/admin/users/{id}/role", "PUT", userController.AdminUpdateRole, true},
		{"/admin/users/{id}", "DELETE", userController.AdminDeleteUser, true},
	}

	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.JwtAuthMiddleware)

	for _, route := range routes {
		if route.Protected {
			if len(route.Path) >= 6 && route.Path[:6] == "/admin" {
				admin.HandleFunc(route.Path[len("/admin"):], route.Handler).Methods(route.Method)
			} else {
				r.Handle(route.Path, middleware.JwtAuthMiddleware(http.HandlerFunc(route.Handler))).Methods(route.Method)
			}
		} else {
			r.HandleFunc(route.Path, route.Handler).Methods(route.Method)
		}
	}

	return r
}
