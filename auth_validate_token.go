package main

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
)


func validateToken(tokenString string) (bool, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Validate the signing method of the token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Provide the secret key used to sign the token
		// Replace "your-secret-key" with your actual secret key
		secretKey := []byte("your-secret-key")
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}

	// Check if the token is valid
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}

	return false, nil
}