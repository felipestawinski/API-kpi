package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/felipestawinski/API-kpi/models"
	"github.com/felipestawinski/API-kpi/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	db := mongoClient
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Resolve the user's ObjectID from their username
	usersCollection := db.Database(database.DbName).Collection(database.CollectionName)
	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err := usersCollection.FindOne(ctx, bson.M{"username": request.Username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Query the files collection for this user's files
	filesCollection := db.Database(database.DbName).Collection(database.FilesCollectionName)
	cursor, err := filesCollection.Find(ctx, bson.M{"ownerId": user.ID})
	if err != nil {
		http.Error(w, "Error finding files", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var files []models.File
	if err := cursor.All(ctx, &files); err != nil {
		http.Error(w, "Error decoding files", http.StatusInternalServerError)
		return
	}

	if files == nil {
		files = []models.File{}
	}

	fmt.Printf("files: %v\n", files)

	// Return the list of files as proper JSON objects
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}