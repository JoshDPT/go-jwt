package main

import (
	"encoding/json"
	"net/http"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

)

func getUsers (w http.ResponseWriter, r *http.Request) {
	
	users, err := getAllUsers()
	if err != nil {
			// Handle the error
			fmt.Println("error getting users")
	}

	// Return the users as JSON response
	json.NewEncoder(w).Encode(users)
	fmt.Println("Serving all users")
}

func getAllUsers() ([]User, error) {
	query := "SELECT id, username, password FROM users"

	rows, err := db.Query(query)
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Username, &user.Password)
			if err != nil {
					return nil, err
			}
			users = append(users, user)
	}

	if err := rows.Err(); err != nil {
			return nil, err
	}

	return users, nil
}







