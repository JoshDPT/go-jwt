package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// getMyEnv retrieves the value of the environment variable specified by the input string.
// It uses the godotenv package to load environment variables from the .env file.
// If the .env file cannot be loaded or the specified environment variable is not found,
// it logs a fatal error and exits the program.
// The function returns the value of the environment variable as a string.
func getMyEnv(s string) string {
	// Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the value of the specified environment variable
	secret := os.Getenv(s)

	return secret
}
