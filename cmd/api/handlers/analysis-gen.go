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
    "fmt"
    "bytes"
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

    fmt.Print("Request body: ", r.Body)

	// Parse the body to get the file ID
	var request struct {
		FileID int `json:"fileId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        fmt.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

    fmt.Println("Request received for file ID:", request.FileID)

	var user models.User
    err = collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Search for the file with the specified ID in the user's files
    var targetFile models.File
    var fileFound bool

    for _, file := range user.Files {
        if file.ID == request.FileID {
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

    // Prepare the JSON payload
    payload := map[string]interface{}{
        "fileAddress": targetFile.FileAddress,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        http.Error(w, "Failed to encode analysis request payload", http.StatusInternalServerError)
        return
    }

	analysisReq, err := http.NewRequest("POST", "http://localhost:9090/analysis-gen", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Failed to create analysis request", http.StatusInternalServerError)
		return
	}
	analysisReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	analysisResp, err := client.Do(analysisReq)
	if err != nil {
		http.Error(w, "Failed to send analysis request", http.StatusInternalServerError)
		return
	}
	defer analysisResp.Body.Close()

	var result map[string]interface{}
    if err := json.NewDecoder(analysisResp.Body).Decode(&result); err != nil {
        http.Error(w, "Failed to decode analysis response", http.StatusInternalServerError)
        return
    }

    fmt.Println("Analysis response:", result)

    // Prepare response data
    responseData := map[string]interface{}{
        "filename": targetFile.Filename,
        "id":       targetFile.ID,
    }

    // Check if figure data exists and add it to response
    if fig, exists := result["fig"]; exists {
        if figStr, ok := fig.(string); ok {
            fmt.Printf("Figure data type: %T\n", fig)
            fmt.Printf("Figure data length: %d characters\n", len(figStr))
            
            // Add the base64 image to the response
            responseData["image"] = figStr
            responseData["hasImage"] = true
        } else {
            fmt.Println("Fig value is not a string")
            responseData["hasImage"] = false
            responseData["error"] = "Invalid image data format"
        }
    } else {
        fmt.Println("Key 'fig' not found in analysis response")
        responseData["hasImage"] = false
        
        // Include any message from the analysis service
        if message, exists := result["message"]; exists {
            responseData["message"] = message
        }
    }

    // Return the complete response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(responseData)
}

