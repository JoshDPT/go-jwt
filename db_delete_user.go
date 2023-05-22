package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// deleteUser handles the deletion of a user based on the provided "id" query parameter.
// It first checks if the "id" parameter is provided and returns a 400 Bad Request status if it is missing.
// Then, it attempts to delete the user with the corresponding ID from the database using the deleteUserByID function.
// If there is an error deleting the user, it returns a 500 Internal Server Error status.
// Otherwise, it returns a success response indicating that the user has been deleted.
//
// The success response is encoded as JSON and sent in the response body.

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	// Get the value of the "id" query parameter
	id := r.URL.Query().Get("id")

	// Handle the case where ID is not provided
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	// Delete the user by ID
	err := deleteUserByID(id)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]string{
		"message": fmt.Sprintf("User with ID %s has been deleted", id),
	}
	log.Printf("Delete user ID: %s successful", id)
	json.NewEncoder(w).Encode(response)
}

// deleteUserByID deletes the user from the database with the provided ID.
// It executes a SQL delete query to remove the user from the "users" table.
// If there is an error executing the delete query, it returns the error.
// Otherwise, it returns nil to indicate successful deletion.

func deleteUserByID(id string) error {

	// Initialize the SQL delete query
	query := "DELETE FROM users WHERE id = ?"

	// Execute the delete query
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
