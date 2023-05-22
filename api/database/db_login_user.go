package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/JoshDPT/go-jwt/api/lib"

	"github.com/golang-jwt/jwt"
	_ "github.com/mattn/go-sqlite3"
)

var jwtSecret = []byte(lib.GetMyEnv("JWT_TOKEN"))

// login authenticates a user by validating their credentials against
// the database and generates a JWT token upon successful authentication.
// It expects a JSON payload containing the user's username and password in the request body.
//
// If the request body is not a valid JSON or the user's credentials are
// invalid or not found, it returns a 401 Unauthorized status.
//
// Upon successful authentication, it generates a JWT token with the user's
// ID and username, signed with a secret key, and returns the signed token in the response.
//
// The token has an expiration time set to 24 hours from the current time.
//
// If there is an error during the authentication or token generation process,
// it returns a 500 Internal Server Error status.

func Login(w http.ResponseWriter, r *http.Request) {
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
	log.Println("User authentication successful")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": signedToken})
}
