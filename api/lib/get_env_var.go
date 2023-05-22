package lib

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

// getMyEnv retrieves the value of the environment variable specified by the input string.
// It uses the godotenv package to load environment variables from the .env file.
// If the .env file cannot be loaded or the specified environment variable is not found,
// it logs a fatal error and exits the program.
// The function returns the value of the environment variable as a string.

func GetMyEnv(s string) string {

	// Get the absolute file path of the root directory
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Construct the absolute file path of the .env file
	envFile := filepath.Join(rootDir, ".env")

	// Load the .env file
	err = godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the value of the specified environment variable
	secret := os.Getenv(s)

	return secret
}
