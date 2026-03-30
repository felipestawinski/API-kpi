package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/felipestawinski/API-kpi/models"
	"github.com/felipestawinski/API-kpi/pkg/config"
	"github.com/felipestawinski/API-kpi/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
)

type filePreviewRequest struct {
	FileID       int  `json:"fileId"`
	MaxRows      int  `json:"maxRows"`
	MaxCols      int  `json:"maxCols"`
	ForceRefresh bool `json:"forceRefresh"`
}

func FilePreviewHandler(w http.ResponseWriter, r *http.Request) {
	if !UserAuthorized(w, r, models.UserStatus(0)) {
		return
	}

	tokenStr := r.Header.Get("Authorization")
	username, err := getUsernameFromToken(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request filePreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.FileID == 0 {
		http.Error(w, "fileId is required", http.StatusBadRequest)
		return
	}

	if request.MaxRows <= 0 {
		request.MaxRows = 20
	}
	if request.MaxCols <= 0 {
		request.MaxCols = 12
	}

	db := database.NewMongoDB(config.MongoURI)
	collection := db.Database(database.DbName).Collection(database.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err = collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var selectedFile *models.File
	for index := range user.Files {
		if user.Files[index].ID == request.FileID {
			selectedFile = &user.Files[index]
			break
		}
	}

	if selectedFile == nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	payload := map[string]interface{}{
		"fileAddress":  selectedFile.FileAddress,
		"fileType":     selectedFile.FileType,
		"maxRows":      request.MaxRows,
		"maxCols":      request.MaxCols,
		"forceRefresh": request.ForceRefresh,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to encode preview request payload", http.StatusInternalServerError)
		return
	}

	previewReq, err := http.NewRequest("POST", "http://127.0.0.1:9090/preview-gen", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Failed to create preview request", http.StatusInternalServerError)
		return
	}
	previewReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	previewResp, err := client.Do(previewReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send preview request: %v", err), http.StatusInternalServerError)
		return
	}
	defer previewResp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(previewResp.StatusCode)
	_, _ = io.Copy(w, previewResp.Body)
}
