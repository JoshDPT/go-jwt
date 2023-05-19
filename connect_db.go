package main

import (
	"log"
	"fmt"
	"database/sql"
)

func connectDB() {
	// Open the SQLite database file
	var err error
	db, err = sql.Open("sqlite3", "login-test.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Opened database file")

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