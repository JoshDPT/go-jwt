package main

import (
	"database/sql"
	"log"
)

func connectTestDB() {
	// Open the SQLite database file
	var err error
	db, err = sql.Open("sqlite3", "db_user_test.db")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database")

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

	// Check if the user with a specific username exists
	username := "josh"
	query := "SELECT COUNT(*) FROM users WHERE username = ?"
	var count int
	err = db.QueryRow(query, username).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	// If the user doesn't exist, insert a new row
	if count == 0 {
		insertUser := "INSERT INTO users (username, password) VALUES (?, ?)"
		_, err = db.Exec(insertUser, username, "josh")
		if err != nil {
			log.Fatal(err)
		}
	}

	// This is for testing purposes to create JWT token

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

