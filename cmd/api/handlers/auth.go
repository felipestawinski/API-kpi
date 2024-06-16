package handlers

import (
	"fmt"
	"net/http"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"github.com/BloxBerg-UTFPR/API-Blockchain/pkg/database"
	"github.com/BloxBerg-UTFPR/API-Blockchain/pkg/config"
	"github.com/BloxBerg-UTFPR/API-Blockchain/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)



// welcomeHandler handles welcome requests for logged-in users
func UserAuthorized(w http.ResponseWriter, r *http.Request) {

	//Check for the Authorization header
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	client := database.NewMongoDB(config.MongoURI)

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		fmt.Println("TOKEN:", token)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	collection := client.Database(database.DbName).Collection(database.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find user by username
	var user models.User
	err = collection.FindOne(ctx, bson.M{"username": claims.Subject}).Decode(&user)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check permission level
	if user.Permission <= 2 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	return
}