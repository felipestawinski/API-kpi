package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
	"fmt"	
    "github.com/felipestawinski/API-kpi/models"
    "github.com/felipestawinski/API-kpi/pkg/config"
    "github.com/felipestawinski/API-kpi/pkg/database"
    "go.mongodb.org/mongo-driver/bson"
)

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the username from the request body
    var request struct {
        Username string `json:"username"`
    }
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Connect to the database
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

	fmt.Printf("files: %v\n", user.Files)

    // Return the list of files in JSON format
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(user.Files); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}