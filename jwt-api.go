package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

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

	router := mux.NewRouter()

	// API endpoint for user authentication
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/users", createUser).Methods("POST")


	log.Fatal(http.ListenAndServe(":8000", router))
}

func connectDB() {
	// Open the SQLite database file
	var err error
	db, err = sql.Open("sqlite3", "login-test.db")
	if err != nil {
		log.Fatal(err)
	}

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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": signedToken})
}
