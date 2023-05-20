package main

import (
	"database/sql"
	"log"
	"net/http"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var db *sql.DB
var jwtSecret = []byte(getMyEnv("JWT_TOKEN")) // Replace with your own secret key - from env file

func main() {
	// Connect to the database
	connectDB()

	mainRouter := mux.NewRouter()

	// Unprotected login route to get JWT token
	subRouter1 := mainRouter.PathPrefix("/api/v1").Subrouter()
	subRouter1.HandleFunc("/login", login).Methods("POST")
	
	// JWT protected routers
	subRouter2 := mainRouter.PathPrefix("/api/v2").Subrouter()

	// AUTH middleware
	// subRouter2.Use(authMiddleware)

	// JWT protected subrouter endpoints
	subRouter2.HandleFunc("/users", createUser).Methods("POST")
	subRouter2.HandleFunc("/users", getUsers).Methods("GET")
	subRouter2.HandleFunc("/users", updateUser).Methods("PUT")
	subRouter2.HandleFunc("/users", deleteUser).Methods("DELETE")

	fmt.Println("Listening on Port 8000")
	log.Fatal(http.ListenAndServe(":8000", mainRouter))
}