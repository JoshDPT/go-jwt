package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"log"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Test_login(t *testing.T) {
	connectTestDB()
	// Set up a test HTTP server
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

	jwtToken := response.Token
	log.Println(jwtToken)
	// Now you can use the `jwtToken` variable to assert or perform further actions based on the received JWT token.
	// For example, you can log it, save it, or use it for subsequent authenticated requests.


	server2 := httptest.NewServer(http.HandlerFunc(createUser))
	defer server2.Close()

	// Define the request body as a byte slice
	requestBody2 := []byte(`{"username": "jimmy", "password": "john"}`)

	// Create a new POST request with the request body
	req2, err := http.NewRequest(http.MethodPost, server2.URL+"/api/v2/users", bytes.NewBuffer(requestBody2))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req2.Header.Add("Content-Type", "application/json")
	req2.Header.Add("Authorization", "Bearer "+jwtToken)
	

	// Send the request
	client2 := &http.Client{}
	resp2, err := client2.Do(req2)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp2.Body.Close()

	// Check the response status code
	if resp2.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 200, but got %d", resp2.StatusCode)
	}




	server3 := httptest.NewServer(http.HandlerFunc(getUsers))
	defer server3.Close()

	// Create a new GET request with the request body
	req3, err := http.NewRequest(http.MethodGet, server3.URL+"/api/v2/users", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req3.Header.Add("Content-Type", "application/json")
	req3.Header.Add("Authorization", "Bearer "+jwtToken)
	


	// Send the request
	client3 := &http.Client{}
	resp3, err := client3.Do(req3)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp3.Body.Close()

	// Read the response body
	body3, err := ioutil.ReadAll(resp3.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Define a slice to store the objects
	var users []User

	// Unmarshal the JSON array into the slice of objects
	err = json.Unmarshal(body3, &users)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	var id int

	log.Println(users)
	
	for _,o := range users {
		if o.Username == "jimmy" {
			id = o.ID
		}
	}

	log.Println(id)
	// Check the response status code
	if resp3.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", resp3.StatusCode)
	}




	server4 := httptest.NewServer(http.HandlerFunc(deleteUser))
	defer server4.Close()

	// Create a new POST request with the request body
	req4, err := http.NewRequest(http.MethodDelete, server4.URL+"/api/v2/users?id="+fmt.Sprint(id), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set the request headers if necessary
	req4.Header.Set("Content-Type", "application/json")

	// Set the authorization header with the JWT token
	req4.Header.Set("Authorization", "Bearer "+jwtToken)


	// Send the request
	client4 := &http.Client{}
	resp4, err := client4.Do(req4)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp4.Body.Close()

	// Check the response status code
	if resp4.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", resp4.StatusCode)
	}
}
