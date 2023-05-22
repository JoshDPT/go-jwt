package main

import (
	"fmt"
	"strings"


	"github.com/golang-jwt/jwt"
	_ "github.com/mattn/go-sqlite3"
)

var jwtSecret = []byte(GetMyEnv("JWT_TOKEN"))
// validateToken validates the provided JWT token string.
// It removes the "Bearer " prefix from the token string, if present.
// Then, it parses the token and validates its signature using the secret key.
// If the token is valid, it extracts the required information from the claims.
// Additional validation logic can be implemented if needed.
// It returns a boolean indicating whether the token is valid and an error, if any.
// If there is an error parsing or validating the token, it logs the error and returns false.

func validateToken(tokenString string) (bool, error) {
	// fmt.Println("inside Validate Token Func", tokenString)

	// Remove the "Bearer " prefix from the token string
	if len(tokenString) > 7 && strings.ToUpper(tokenString[0:7]) == "BEARER " {
		tokenString = tokenString[7:]
	}

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method of the token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// Check if the token is valid and has a valid signature
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
