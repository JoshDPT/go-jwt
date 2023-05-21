package main

import (
	"database/sql"
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

var db *sql.DB
var jwtSecret = []byte(getMyEnv("JWT_TOKEN"))

func main() {
	// Connect to the database
	connectDB()

	mainRouter := mux.NewRouter()

	// Unprotected login route to get JWT token
	subRouter1 := mainRouter.PathPrefix("/api/v1").Subrouter()
	subRouter1.HandleFunc("/login", login).Methods("POST")

	// JWT protected routers
	subRouter2 := mainRouter.PathPrefix("/api/v2").Subrouter()

	// Use AUTH middleware
	subRouter2.Use(authMiddleware)

	// Define JWT protected subrouter endpoints
	subRouter2.HandleFunc("/users", createUser).Methods("POST")
	subRouter2.HandleFunc("/users", getUsers).Methods("GET")
	subRouter2.HandleFunc("/users", updateUser).Methods("PUT")
	subRouter2.HandleFunc("/users", deleteUser).Methods("DELETE")

	log.Println("Server is available at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", mainRouter))
}

