package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/felipestawinski/API-kpi/models"
	"github.com/felipestawinski/API-kpi/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// Do not create a rigid short-lived context here since the Python analysis takes a long time.
	db := mongoClient
	collection := db.Database(database.DbName).Collection(database.CollectionName)
	
	initCtx, initCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer initCancel()

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
	err = collection.FindOne(initCtx, bson.M{"username": username}).Decode(&user)
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

	// ── Streaming text path ──────────────────────────────────────────────────
	// When the Python service streams plain text (text-only analysis), forward
	// the body directly to the client without JSON decoding or token tracking.
	contentType := analysisResp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/plain") {
		fmt.Println("Forwarding streaming text response from Python service")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		if _, copyErr := io.Copy(w, analysisResp.Body); copyErr != nil {
			fmt.Println("Error forwarding streaming body:", copyErr)
		}
		return
	}

	// ── JSON path (chart / chart recommendation) ──────────────────────────────
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

	// Track token usage from the analysis service response
	tokensConsumed := 0
	if tu, exists := result["tokens_used"]; exists && tu != nil {
		switch v := tu.(type) {
		case float64:
			tokensConsumed = int(v)
		case int:
			tokensConsumed = v
		}
	}

	fmt.Printf("DEBUG_TOKENS: Extracted %d tokens from Python response\n", tokensConsumed)

	var tokenUserResult models.User
	if tokensConsumed > 0 {
		// Atomically increment tokensUsed and return the updated document in one round-trip.
		updateCtx, updateCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer updateCancel()

		after := options.After
		findOneAndUpdateErr := collection.FindOneAndUpdate(
			updateCtx,
			bson.M{"username": username},
			bson.M{"$inc": bson.M{"tokensUsed": tokensConsumed}},
			&options.FindOneAndUpdateOptions{ReturnDocument: &after},
		).Decode(&tokenUserResult)
		if findOneAndUpdateErr != nil {
			fmt.Println("Warning: failed to update token usage:", findOneAndUpdateErr)
			// Fall back to the user loaded at the start of the handler.
			tokenUserResult = user
		} else {
			fmt.Printf("DEBUG_TOKENS: FindOneAndUpdate succeeded for user %s, tokensUsed now: %d\n", username, tokenUserResult.TokensUsed)
		}
	} else {
		// No tokens consumed – reuse the user struct already in memory.
		tokenUserResult = user
	}

	tokenLimit := tokenUserResult.TokenLimit
	if tokenLimit == 0 {
		tokenLimit = 1000000
	}
	responseData["tokensUsed"] = tokenUserResult.TokensUsed
	responseData["tokenLimit"] = tokenLimit

	// Return the complete response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
