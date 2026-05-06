package handlers

import (
	"bytes"
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

type filePreviewCachedRequest struct {
	FileID string `json:"fileId"`
}

// FilePreviewCachedHandler returns a cached file preview from MongoDB.
// If the file doesn't have a cached preview (uploaded before this feature),
// it falls back to calling the Python /preview-gen endpoint, caches the
// result in MongoDB for future use, and returns the preview.
func FilePreviewCachedHandler(w http.ResponseWriter, r *http.Request) {
	if !UserAuthorized(w, r, models.UserStatus(0)) {
		return
	}

	tokenStr := r.Header.Get("Authorization")
	username, err := getUsernameFromToken(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request filePreviewCachedRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.FileID == "" {
		http.Error(w, "fileId is required", http.StatusBadRequest)
		return
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

	// If preview is cached, return it directly from MongoDB
	if len(selectedFile.PreviewHeaders) > 0 {
		fmt.Printf("FilePreviewCached: returning cached preview for file %s (%d headers, %d rows)\n",
			request.FileID, len(selectedFile.PreviewHeaders), len(selectedFile.PreviewRows))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"headers": selectedFile.PreviewHeaders,
			"rows":    selectedFile.PreviewRows,
		})
		return
	}

	// Fallback: file was uploaded before this feature — fetch from Python
	fmt.Printf("FilePreviewCached: no cached preview for file %s, falling back to Python\n", request.FileID)

	payload := map[string]interface{}{
		"fileAddress":  selectedFile.FileAddress,
		"fileType":     selectedFile.FileType,
		"maxRows":      200,
		"maxCols":      100,
		"forceRefresh": false,
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

	client := &http.Client{Timeout: 30 * time.Second}
	previewResp, err := client.Do(previewReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch preview: %v", err), http.StatusInternalServerError)
		return
	}
	defer previewResp.Body.Close()

	if previewResp.StatusCode != http.StatusOK {
		http.Error(w, "Preview service returned an error", previewResp.StatusCode)
		return
	}

	var result struct {
		Headers []string   `json:"headers"`
		Rows    [][]string `json:"rows"`
	}
	if err := json.NewDecoder(previewResp.Body).Decode(&result); err != nil {
		http.Error(w, "Failed to decode preview response", http.StatusInternalServerError)
		return
	}

	// Cache the preview in the files collection for future requests (fire-and-forget)
	go func() {
		cacheCtx, cacheCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cacheCancel()

		filter := bson.M{"_id": fileObjID}
		update := bson.M{
			"$set": bson.M{
				"previewHeaders": result.Headers,
				"previewRows":    result.Rows,
			},
		}
		_, updateErr := filesCollection.UpdateOne(cacheCtx, filter, update)
		if updateErr != nil {
			fmt.Printf("FilePreviewCached: failed to cache preview for file %s: %v\n", request.FileID, updateErr)
		} else {
			fmt.Printf("FilePreviewCached: cached preview for file %s (%d headers, %d rows)\n",
				request.FileID, len(result.Headers), len(result.Rows))
		}
	}()

	// Return the preview
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"headers": result.Headers,
		"rows":    result.Rows,
	})
}
