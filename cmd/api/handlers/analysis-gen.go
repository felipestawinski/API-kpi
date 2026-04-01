package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/felipestawinski/API-kpi/models"
	"github.com/felipestawinski/API-kpi/pkg/config"
	"github.com/felipestawinski/API-kpi/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net/http"
	"strings"
	"time"
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

	// Parse the body to get the file IDs (now accepting multiple)
	var request struct {
		FileIDs             []int  `json:"fileIds"`
		Prompt              string `json:"prompt"`
		GenerateChart       bool   `json:"generateChart"`
		ChartRecommendation bool   `json:"chartRecommendation"`
		ChatID              string `json:"chatId"`
		ForceRefresh        bool   `json:"forceRefresh"`
		Model               string `json:"model"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate that at least one file ID was provided
	if len(request.FileIDs) == 0 {
		http.Error(w, "At least one file ID is required", http.StatusBadRequest)
		return
	}

	request.Model = strings.TrimSpace(request.Model)
	if request.Model == "" {
		request.Model = "gpt-5-mini"
	}

	var user models.User
	err = collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Search for all files with the specified IDs and collect their addresses
	fileAddresses := []string{}
	fileTypes := []string{}
	foundFileIDs := []int{}

	for _, requestedID := range request.FileIDs {
		for _, file := range user.Files {
			if file.ID == requestedID {
				fileAddresses = append(fileAddresses, file.FileAddress)
				fileTypes = append(fileTypes, file.FileType)
				foundFileIDs = append(foundFileIDs, file.ID)
				break
			}
		}
	}

	// Check if all requested files were found
	if len(foundFileIDs) != len(request.FileIDs) {
		http.Error(w, fmt.Sprintf("Not all files found. Requested: %d, Found: %d", len(request.FileIDs), len(foundFileIDs)), http.StatusNotFound)
		return
	}

	fmt.Println("Generating analysis for files with IDs:", foundFileIDs)
	fmt.Println("File addresses:", fileAddresses)

	fmt.Println("File addresses to be sent for analysis:", fileAddresses)
	// Prepare the JSON payload with multiple file addresses
	payload := map[string]interface{}{
		"fileAddresses":       fileAddresses,
		"fileTypes":           fileTypes,
		"prompt":              request.Prompt,
		"generateChart":       request.GenerateChart,
		"chartRecommendation": request.ChartRecommendation,
		"chatId":              request.ChatID,
		"forceRefresh":        request.ForceRefresh,
		"model":               request.Model,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to encode analysis request payload:", err)
		http.Error(w, "Failed to encode analysis request payload", http.StatusInternalServerError)
		return
	}

	analysisReq, err := http.NewRequest("POST", "http://127.0.0.1:9090/analysis-gen", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Failed to create analysis request:", err)
		http.Error(w, "Failed to create analysis request", http.StatusInternalServerError)
		return
	}
	analysisReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 150 * time.Second}
	analysisResp, err := client.Do(analysisReq)
	if err != nil {
		fmt.Println("Failed to send analysis request:", err)
		http.Error(w, "Failed to send analysis request", http.StatusInternalServerError)
		return
	}
	defer analysisResp.Body.Close()

	if analysisResp.StatusCode < 200 || analysisResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(analysisResp.Body)
		if len(bodyBytes) == 0 {
			http.Error(w, "Analysis service returned an error", analysisResp.StatusCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(analysisResp.StatusCode)
		_, _ = w.Write(bodyBytes)
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(analysisResp.Body).Decode(&result); err != nil {
		fmt.Println("Failed to decode analysis response:", err)
		http.Error(w, "Failed to decode analysis response", http.StatusInternalServerError)
		return
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"ids": foundFileIDs,
	}

	// Check if figure data exists and is a non-null string
	if fig, exists := result["chart_base64"]; exists && fig != nil {
		if figStr, ok := fig.(string); ok && figStr != "" {
			fmt.Printf("Figure data type: %T\n", fig)
			fmt.Printf("Figure data length: %d characters\n", len(figStr))

			// Add the base64 image to the response
			responseData["image"] = figStr
			responseData["hasImage"] = true
		} else {
			fmt.Println("Fig value is not a valid non-empty string")
			responseData["hasImage"] = false
		}
	} else {
		fmt.Println("Key 'chart_base64' not found or is null in analysis response")
		responseData["hasImage"] = false

		// Include any message from the analysis service
		if message, exists := result["message"]; exists {
			responseData["message"] = message
		}
	}

	if text_response, exists := result["text_response"]; exists {
		if textStr, ok := text_response.(string); ok {
			responseData["text_response"] = textStr
		}
	}

	// Return the complete response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
