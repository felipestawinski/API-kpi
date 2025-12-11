package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/felipestawinski/API-kpi/models"
	"github.com/felipestawinski/API-kpi/pkg/config"
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

	// Connect to the database
	db := database.NewMongoDB(config.MongoURI)
	collection := db.Database(database.DbName).Collection(database.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find all users
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error finding users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var allFiles []models.File
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			fmt.Printf("Error decoding user: %v\n", err)
			continue
		}

		// Parse files for this user based on searchType
		for _, file := range user.Files {
			var match bool
			switch request.SearchType {
			case "filename":
				match = file.Filename == request.SearchTerm
			case "institution":
				match = file.Institution == request.SearchTerm
			case "writer":
				match = file.Writer == request.SearchTerm
			}
			if match {
				allFiles = append(allFiles, file)
			}
		}
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
