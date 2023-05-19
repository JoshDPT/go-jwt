package main

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
)

func validateToken(tokenString string) (bool, error) {
	fmt.Println("inside Validate Token Func", tokenString)

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method of the token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Provide the secret key used to sign the token
		// Replace "your-secret-key" with your actual secret key
		secretKey := []byte("1234")
		return secretKey, nil
	})

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// Check if the token is valid and has valid signature
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// Extract the required information from claims, if needed
		userID := claims["id"].(float64)
		username := claims["username"].(string)
		fmt.Println(userID, username)

		// Additional validation logic, if required

		return true, nil
	}

	fmt.Println("token not valid")
	return false, nil
}
