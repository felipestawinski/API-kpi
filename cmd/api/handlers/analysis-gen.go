package handlers

import (
	"fmt"
	"net/http"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"github.com/felipestawinski/API-kpi/pkg/database"
	"github.com/felipestawinski/API-kpi/pkg/config"
	"github.com/felipestawinski/API-kpi/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func AnalysisGenHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	username, err := getUsernameFromToken(tokenStr)

	// Check if the user is authorized
    if !UserAuthorized(w, r, models.UserStatus(0)) {
        return 
    }

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db := database.NewMongoDB(config.MongoURI)
	collection := db.Database(database.DbName).Collection(database.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Parse the body to get the file ID
	var request stru


}