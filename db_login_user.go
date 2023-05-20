package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
)

func login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Query the database for the user's credentials
	query := "SELECT id, username, password FROM users WHERE username = ?"
	row := db.QueryRow(query, user.Username)

	var dbUser User
	err = row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Validate the password
	if user.Password != dbUser.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       dbUser.ID,
		"username": dbUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (e.g., 24 hours from now)
	})

	// Sign the token with your secret key
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the signed JWT token
	fmt.Println("Success in authenticating user")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": signedToken})
}