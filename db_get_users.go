package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	// Get the value of the "id" query parameter
	id := r.URL.Query().Get("id")

	if id != "" {
		// Handle a specific user request based on the ID
		user, err := getUserByID(id)
		if err != nil {
			// Handle the error
			http.Error(w, "Error getting user by ID", http.StatusInternalServerError)
			return
		}

		// Return the user as JSON response
		json.NewEncoder(w).Encode(user)
		log.Printf("Serve user ID: %s successful", id)
	} else {
		// Handle the request to get all users
		users, err := getAllUsers()
		if err != nil {
			// Handle the error
			http.Error(w, "Error getting users", http.StatusInternalServerError)
			return
		}

		// Return the users as JSON response
		json.NewEncoder(w).Encode(users)
		log.Println("Serve all users successful")
	}
}

// Func for getting all users is no query params
func getAllUsers() ([]User, error) {

	// Initialize SQL query as variable
	query := "SELECT id, username, password FROM users"

	// Query the database with the query string
	rows, err := db.Query(query)

	// Handle the error
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// initialize a slice of Users
	var users []User

	// iterate through the matching rows and add it to the slice
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Handle the error
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Use this func if have id query params
func getUserByID(id string) (User, error) {

	// Initialize query string for specific id
	query := "SELECT id, username, password FROM users WHERE id = ?"

	// Query db with string & id
	row := db.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password)

	// Handle error with scanning
	if err != nil {
		return User{}, err
	}

	return user, nil
}
