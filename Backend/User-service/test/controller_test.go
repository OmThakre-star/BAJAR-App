package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/controller"
)

func TestRegister(t *testing.T) {
	req := httptest.NewRequest("POST", "/register", nil)
	w := httptest.NewRecorder()
	controller.Register(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
		t.Errorf("unexpected status code: %d", w.Code)
	}
}

func TestLogin(t *testing.T) {
	req := httptest.NewRequest("POST", "/login", nil)
	w := httptest.NewRecorder()
	controller.Login(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized && w.Code != http.StatusBadRequest {
		t.Errorf("unexpected status code: %d", w.Code)
	}
}

func TestGetProfile(t *testing.T) {
	req := httptest.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, "userID", 1) // Set dummy userID
	req = req.WithContext(ctx)
	controller.GetProfile(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("unexpected status code: %d", w.Code)
	}
}

func TestListUsers(t *testing.T) {
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	controller.ListUsers(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code: %d", w.Code)
	}
}
