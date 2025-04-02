package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/BloxBerg-UTFPR/API-Blockchain/models"
    "github.com/BloxBerg-UTFPR/API-Blockchain/pkg/config"
    "github.com/BloxBerg-UTFPR/API-Blockchain/pkg/database"
    "go.mongodb.org/mongo-driver/bson"
	"time"
)

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {

	UserAuthorized(w, r)
	
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
	    "accessType":  user.Permission,
		"position"	: user.Role,
	    "accessTime": "0",})

	return
}