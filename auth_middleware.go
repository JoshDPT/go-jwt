package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// authMiddleware is a middleware function that performs authentication.
// It expects an Authorization header with a valid JWT token.
// It validates the token using the `validateToken` function.
// If the token is not provided or not valid, it returns a 401 Unauthorized status.
// If the token is valid, it allows the request to proceed to the next handler.
// It logs whether the token is valid or not.
// It takes the next http.Handler as input and returns an http.Handler.


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
			log.Println("Token is not valid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next handler
		log.Println("Token is valid")
		next.ServeHTTP(w, r)
	})
}
