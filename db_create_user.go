package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// createUser handles the creation of a new user based on the provided request body.
// It decodes the JSON request body into a User struct and performs validation on the user data.
// If the user data is invalid (empty username or password), it returns a 400 Bad Request status.
// Otherwise, it inserts the user into the database using a prepared statement.
// If there is an error inserting the user, it returns a 500 Internal Server Error status.
// Otherwise, it returns a success response with a 201 Created status.

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Perform validation on the user data
	if user.Username == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Insert the user into the database
	stmt, err := db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(user.Username, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return a success response
	log.Println("User creation successful")
	w.WriteHeader(http.StatusCreated)
}
