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
	"github.com/felipestawinski/API-kpi/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type filePreviewRequest struct {
	FileID       string `json:"fileId"`
	MaxRows      int    `json:"maxRows"`
	MaxCols      int    `json:"maxCols"`
	ForceRefresh bool   `json:"forceRefresh"`
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

	if request.FileID == "" {
		http.Error(w, "fileId is required", http.StatusBadRequest)
		return
	}

	if request.MaxRows <= 0 {
		request.MaxRows = 20
	}
	if request.MaxCols <= 0 {
		request.MaxCols = 12
	}

	fileObjID, err := primitive.ObjectIDFromHex(request.FileID)
	if err != nil {
		http.Error(w, "Invalid fileId format", http.StatusBadRequest)
		return
	}

	db := mongoClient
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Resolve user ObjectID
	usersCollection := db.Database(database.DbName).Collection(database.CollectionName)
	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err = usersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Find the file directly in the files collection
	filesCollection := db.Database(database.DbName).Collection(database.FilesCollectionName)
	var selectedFile models.File
	err = filesCollection.FindOne(ctx, bson.M{"_id": fileObjID, "ownerId": user.ID}).Decode(&selectedFile)
	if err != nil {
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
