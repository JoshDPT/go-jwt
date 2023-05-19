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
var jwtSecret = []byte("1234") // Replace with your own secret key

func main() {
	// Connect to the database
	connectDB()

	mainRouter := mux.NewRouter()

	// Unprotected login route to get JWT token
	subRouter1 := mainRouter.PathPrefix("/api/v1").Subrouter()
	subRouter1.HandleFunc("/login", login).Methods("POST")
	
	// JWT protected routers
	subRouter2 := mainRouter.PathPrefix("/api/v2").Subrouter()
	subRouter2.Use(authMiddleware)
	subRouter2.HandleFunc("/users", createUser).Methods("POST")
	subRouter2.HandleFunc("/users", getUsers).Methods("GET")

	fmt.Println("Listening on Port 8000")
	log.Fatal(http.ListenAndServe(":8000", mainRouter))
}