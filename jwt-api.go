package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var db *sql.DB
var jwtSecret = []byte("your-secret-key") // Replace with your own secret key

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

	log.Fatal(http.ListenAndServe(":8000", mainRouter))
}

func connectDB() {
	// Open the SQLite database file
	var err error
	db, err = sql.Open("sqlite3", "login-test.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Opened database file")

	// Create the users table if it doesn't exist
	createTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			password TEXT
		)
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("Success in adding user")
	w.WriteHeader(http.StatusCreated)
}

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
	fmt.Println("Success in logging on, returned JWT")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": signedToken})
}


func validateToken(tokenString string) (bool, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Validate the signing method of the token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Provide the secret key used to sign the token
		// Replace "your-secret-key" with your actual secret key
		secretKey := []byte("your-secret-key")
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}

	// Check if the token is valid
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}

	return false, nil
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the authorization header from the request
		tokenString := r.Header.Get("Authorization")

		// Check if the token is provided
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Validate the token
		isValid, err := validateToken(tokenString)
		if err != nil || !isValid {
			fmt.Println("Token is not valid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next handler
		fmt.Println("Token is valid")
		next.ServeHTTP(w, r)
	})
}

