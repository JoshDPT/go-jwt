package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func TestLoginAndCreateUser(t *testing.T) {
	// Set up test environment
	connectTestDB()

	// Perform login and retrieve JWT token
	jwtToken := performLogin_test(t)

	// Use JWT token to create a new user
	createUser_test(t, jwtToken)


	// Use JWT token to get all users and find ID of created user
	id := getUsersAndFindUserID_test(t, jwtToken)

	// se JWT token to delete created user with ID query params
	deleteUserByID_test(t, jwtToken, id)
}


func performLogin_test(t *testing.T) string {

	server := httptest.NewServer(http.HandlerFunc(login))
	defer server.Close()

	// Define the request body as a byte slice
	requestBody := []byte(`{"username": "josh", "password": "josh"}`)

	// Create a new POST request with the request body
	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v1/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set the request headers if necessary
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Extract the JWT token from the response body
	var response struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	return response.Token
}

func createUser_test(t *testing.T, jwtToken string) {
	// Set up a test HTTP server for creating user
	server := httptest.NewServer(http.HandlerFunc(createUser))
	defer server.Close()

	// Define the request body as a byte slice
	requestBody := []byte(`{"username": "jimmy", "password": "john"}`)

	// Create a new POST request with the request body
	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v2/users", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, but got %d", resp.StatusCode)
	}
}

func getUsersAndFindUserID_test(t *testing.T, jwtToken string) int {
	// Set up a test HTTP server for getting users
	server := httptest.NewServer(http.HandlerFunc(getUsers))
	defer server.Close()

	// Create a new GET request
	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/v2/users", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Define a slice to store the objects
	var users []User

	// Unmarshal the JSON array into the slice of objects
	err = json.Unmarshal(body, &users)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	var id int

	for _, o := range users {
		if o.Username == "jimmy" {
			id = o.ID
		}
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", resp.StatusCode)
	}

	return id
}

func deleteUserByID_test(t *testing.T, jwtToken string, id int) {
	// Set up a test HTTP server for deleting user
	server := httptest.NewServer(http.HandlerFunc(deleteUser))
	defer server.Close()

	// Create a new DELETE request
	req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/v2/users?id="+fmt.Sprint(id), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set the request headers if necessary
	req.Header.Set("Content-Type", "application/json")

	// Set the authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", resp.StatusCode)
	}
}
