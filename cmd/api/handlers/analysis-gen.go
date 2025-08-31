package handlers

import (
	"net/http"
	"time"
	"github.com/felipestawinski/API-kpi/pkg/database"
	"github.com/felipestawinski/API-kpi/pkg/config"
	"github.com/felipestawinski/API-kpi/models"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"fmt"
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
	var request struct {
		FileID string `json:"file_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        fmt.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user models.User
    err = collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Search for the file with the specified ID in the user's files
    var targetFile models.File
    var fileFound bool

    fileID, err := strconv.Atoi(request.FileID)
    if err != nil {
        http.Error(w, "Invalid file ID format", http.StatusBadRequest)
        return
    }

    for _, file := range user.Files {
        if file.ID == fileID {
            targetFile = file
            fileFound = true
            break
        }
    }

    if !fileFound {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    fmt.Println("Generating analysis for file:", targetFile.Filename)
    fmt.Println("File address:", targetFile.FileAddress)

    // Return the file address
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "fileAddress": targetFile.FileAddress,
        "filename":    targetFile.Filename,
        "id":          targetFile.ID,
    })
}

