package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/felipestawinski/API-kpi/models"
    "github.com/felipestawinski/API-kpi/pkg/config"
    "github.com/felipestawinski/API-kpi/pkg/database"
    "go.mongodb.org/mongo-driver/bson"
	"time"
)

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {

	if !UserAuthorized(w, r, models.UserStatus(1)) {
		return
	}
	
	var request struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("request", request)

	//Connect to the database
	db := database.NewMongoDB(config.MongoURI)
	collection := db.Database(database.DbName).Collection(database.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find the user by username
	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": request.Username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"username":     user.Username,
		"institution":  user.Institution,
		"email":        user.Email,
	    "accessType":  models.UserStatus(user.Permission).String(),
		"position"	: user.Role,
	    "accessTime": user.AccessTime,})

	return
}