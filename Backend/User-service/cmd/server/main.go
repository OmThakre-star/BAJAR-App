package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"user-service/config"
	"user-service/controller"
	"user-service/db"
	"user-service/middleware"
)

func main() {

	cfg := config.Load()
	db.InitDB(cfg.DatabaseDSN)

	r := mux.NewRouter()
	r.Use(middleware.RateLimitMiddleware)

	// public
	r.HandleFunc("/register", controller.Register).Methods("POST")
	r.HandleFunc("/login", controller.Login).Methods("POST")

	// protected
	auth := r.PathPrefix("/api").Subrouter()
	auth.Use(middleware.JwtAuthMiddleware)
	auth.HandleFunc("/me", controller.GetProfile).Methods("GET")
	auth.HandleFunc("/users", controller.ListUsers).Methods("GET")

	log.Printf("Server running on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))

}
