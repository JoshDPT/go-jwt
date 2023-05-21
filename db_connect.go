package main

import (
	"database/sql"
	"log"
)

func connectDB() {
	// Open the SQLite database file
	var err error
	db, err = sql.Open("sqlite3", "db_user.db")
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

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
