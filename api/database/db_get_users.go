package database

import (
	"encoding/json"
	"log"
	"net/http"
)

// getUsers retrieves either all users or a specific user based on the "id" query parameter.
// If the "id" query parameter is provided, it returns the user with the corresponding ID.
// Otherwise, it returns all users.
//
// If there is an error retrieving the user(s) from the database, it returns a 500 Internal Server Error status.
//
// The returned user(s) are encoded as JSON and sent in the response body.

func GetUsers(w http.ResponseWriter, r *http.Request) {
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

// getAllUsers retrieves all users from the database.
//
// If there is an error retrieving the users, it returns a nil slice and the error.
//
// Otherwise, it returns a slice of User objects containing the retrieved users.
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

// getUserByID retrieves a specific user from the database based on the provided ID.
//
// If there is an error retrieving the user, it returns an empty User object and the error.
//
// Otherwise, it returns the User object corresponding to the provided ID.
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
