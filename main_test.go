package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
)

func TestCreateUser(t *testing.T) {
	// Create a new HTTP request for creating a user
	reqBody := []byte(`{"username":"testuser","password":"testpassword"}`)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the API response
	rr := httptest.NewRecorder()

	// Create a test router and serve the request
	router := mux.NewRouter()
	router.HandleFunc("/users", createUser).Methods("POST")
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d but got %d", http.StatusCreated, rr.Code)
	}

	// Parse the response body into a JSON object
	var response struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Validate the response body
	expectedUsername := "testuser"
	if response.Username != expectedUsername {
		t.Errorf("expected username %s but got %s", expectedUsername, response.Username)
	}

	// Initialize the database connection
	var db *sql.DB 

	// Query the database to retrieve the newly created user
	var retrievedUser User
	err = db.QueryRow("SELECT id, username FROM users WHERE username = ?", expectedUsername).Scan(&retrievedUser.ID, &retrievedUser.Username)
	if err != nil {
		t.Fatal(err)
	}

	// Validate the retrieved user
	if retrievedUser.Username != expectedUsername {
		t.Errorf("expected username %s but got %s", expectedUsername, retrievedUser.Username)
	}
}
