package main

import (
	"net/http"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the authorization header from the request
		tokenString := r.Header.Get("Authorization")

		// Check if the token is provided
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Validate the token
		isValid, err := validateToken(tokenString)
		if err != nil || !isValid {
			fmt.Println("Token is not valid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next handler
		fmt.Println("Token is valid")
		next.ServeHTTP(w, r)
	})
}
