package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/middleware"
)

func TestJwtAuthMiddleware(t *testing.T) {
	handler := middleware.JwtAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		t.Errorf("should not be OK without token")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	handler := middleware.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusTooManyRequests {
		t.Errorf("unexpected status code: %d", w.Code)
	}
}
