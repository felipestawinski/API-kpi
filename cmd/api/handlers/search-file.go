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
)

func SearchFilesHandler(w http.ResponseWriter, r *http.Request) {
	//Check jwt key
	if !UserAuthorized(w, r, models.UserStatus(0)) {
		return
	}

	// Parse the search parameters from the request body
	var request struct {
		SearchType string `json:"searchType"`
		SearchTerm string `json:"searchTerm"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("request searchType: ", request.SearchType)
	fmt.Println("request searchTerm: ", request.SearchTerm)

	// Validate searchType
	validSearchTypes := map[string]bool{"filename": true, "institution": true, "writer": true}
	if !validSearchTypes[request.SearchType] {
		http.Error(w, "Invalid searchType. Must be 'filename', 'institution', or 'writer'", http.StatusBadRequest)
		return
	}

	if request.SearchTerm == "" {
		http.Error(w, "searchTerm is required", http.StatusBadRequest)
		return
	}

	// Query the files collection directly with an indexed field
	db := mongoClient
	filesCollection := db.Database(database.DbName).Collection(database.FilesCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{request.SearchType: request.SearchTerm}
	cursor, err := filesCollection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Error searching files", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var allFiles []models.File
	if err := cursor.All(ctx, &allFiles); err != nil {
		http.Error(w, "Error decoding files", http.StatusInternalServerError)
		return
	}

	if len(allFiles) == 0 {
		http.Error(w, "No files found matching the search criteria", http.StatusNotFound)
		return
	} else {
		fmt.Println("Found files:", allFiles)
	}

	// Return all matching files in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]models.File{
		"files": allFiles,
	})
}
