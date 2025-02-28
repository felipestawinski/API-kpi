package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/BloxBerg-UTFPR/API-Blockchain/pkg/database"
	"github.com/BloxBerg-UTFPR/API-Blockchain/pkg/config"
	"github.com/BloxBerg-UTFPR/API-Blockchain/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
)


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	client := database.NewMongoDB(config.MongoURI)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginUser models.User
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Login method")
	fmt.Println("user", loginUser)

	collection := client.Database(database.DbName).Collection(database.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find by email
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": loginUser.Email}).Decode(&user)
	if err != nil {
		fmt.Println("aqui1")
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		fmt.Println("aqui2")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create the JWT token
	expirationTime := time.Now().Add(time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   user.Username,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	// Send the token as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"username": user.Username,
		"institution": "UTFPR",
		"token":      tokenString,
		"permission": user.Permission,
	})

	//fmt.Fprintln(w, "Login successful")
}