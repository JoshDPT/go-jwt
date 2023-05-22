package main

import (
	"database/sql"
	"log"
)

// connectDB establishes a connection to the SQLite database file.
// It opens the database file and creates the users table if it doesn't exist.
// If there is an error opening the database file or creating the table, it logs the error and terminates the program.
// After successfully connecting to the database, it logs a success message.
var db *sql.DB

func ConnectDB() {
	// Open the SQLite database file
	var err error
	db, err = sql.Open("sqlite3", "db_user.db")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to database successful")

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

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
