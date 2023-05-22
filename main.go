package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// Connect to the database
	ConnectDB()

	mainRouter := mux.NewRouter()

	// Rate limiting middleware
	mainRouter.Use(RateLimitMiddleware)

	// Unprotected login route to get JWT token
	subRouter1 := mainRouter.PathPrefix("/api/v1").Subrouter()
	subRouter1.HandleFunc("/login", Login).Methods("POST")

	// JWT protected routers
	subRouter2 := mainRouter.PathPrefix("/api/v2").Subrouter()

	// Use AUTH middleware
	subRouter2.Use(AuthMiddleware)

	// Define JWT protected subrouter endpoints
	subRouter2.HandleFunc("/users", CreateUser).Methods("POST")
	subRouter2.HandleFunc("/users", GetUsers).Methods("GET")
	subRouter2.HandleFunc("/users", UpdateUser).Methods("PUT")
	subRouter2.HandleFunc("/users", DeleteUser).Methods("DELETE")

	log.Println("Server is available at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", mainRouter))
}
