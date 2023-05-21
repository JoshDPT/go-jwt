package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func getMyEnv(s string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secret := (os.Getenv(s))
	return secret
}
