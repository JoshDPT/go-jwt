package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func deleteUser(w http.ResponseWriter, r *http.Request) {

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
	json.NewEncoder(w).Encode(response)
}

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
