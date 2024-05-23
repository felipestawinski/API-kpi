package handlers

import (
	"fmt"
	"net/http"
	"github.com/golang-jwt/jwt/v4"
)



// welcomeHandler handles welcome requests for logged-in users
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		fmt.Print("AQUI")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		fmt.Println("TOKEN:", token)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome, %s!", claims.Subject)
}