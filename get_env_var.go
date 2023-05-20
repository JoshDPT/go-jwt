package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func getMyEnv(s string) string {
	err := godotenv.Load(".env")
	if err != nil {
			log.Fatal("Error loading .env file")
	}
	secret := (os.Getenv(s))	
	return secret
}
