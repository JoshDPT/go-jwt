package database

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// updateUser updates a user's information in the database.
// It expects a JSON payload containing user data in the request body.
// The user's username and password are validated, and then the user's
// information is updated in the database.
//
// If the request body is not a valid JSON or the user data is missing,
// it returns a 400 Bad Request status.
//
// If there is an error updating the user in the database, it returns a
// 500 Internal Server Error status.
//
// If the user update is successful, it returns a 200 OK status.

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	// Update the user in the database
	stmt, err := db.Prepare("UPDATE users SET username = ?, password = ? WHERE id = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(user.Username, user.Password, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return a success response
	fmt.Println("Success in updating user")
	w.WriteHeader(http.StatusOK)
}
