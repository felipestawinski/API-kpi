package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/felipestawinski/API-kpi/models"
	"github.com/felipestawinski/API-kpi/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func TokenUsageHandler(w http.ResponseWriter, r *http.Request) {
	if !UserAuthorized(w, r, models.UserStatus(0)) {
		return
	}

	tokenStr := r.Header.Get("Authorization")
	username, err := getUsernameFromToken(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db := mongoClient
	collection := db.Database(database.DbName).Collection(database.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err = collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		fmt.Println("Token usage: user not found:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// For existing users who registered before the token feature, set a default limit
	tokenLimit := user.TokenLimit
	if tokenLimit == 0 {
		tokenLimit = 1000000
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tokensUsed": user.TokensUsed,
		"tokenLimit": tokenLimit,
	})
}
