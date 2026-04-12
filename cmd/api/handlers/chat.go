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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveChatMessageHandler saves a single chat message to MongoDB.
// POST /chat/save
func SaveChatMessageHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	username, err := getUsernameFromToken(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !UserAuthorized(w, r, models.UserStatus(0)) {
		return
	}

	var request struct {
		ChatID    string `json:"chatId"`
		Type      string `json:"type"`
		Content   string `json:"content"`
		Image     string `json:"image,omitempty"`
		Timestamp string `json:"timestamp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println("Error decoding chat save request:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.ChatID == "" || request.Type == "" || request.Content == "" {
		http.Error(w, "chatId, type, and content are required", http.StatusBadRequest)
		return
	}

	ts, err := time.Parse(time.RFC3339, request.Timestamp)
	if err != nil {
		ts = time.Now()
	}

	msg := models.ChatMessage{
		ChatID:    request.ChatID,
		Username:  username,
		Type:      request.Type,
		Content:   request.Content,
		Image:     request.Image,
		HasImage:  request.Image != "",
		Timestamp: ts,
	}

	db := mongoClient
	collection := db.Database(database.DbName).Collection(database.ChatCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, msg)
	if err != nil {
		fmt.Println("Error saving chat message:", err)
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"_id":     result.InsertedID,
		"success": true,
	})
}

// LoadChatHistoryHandler returns all messages for a chatId + username, excluding image data.
// POST /chat/load
func LoadChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	username, err := getUsernameFromToken(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !UserAuthorized(w, r, models.UserStatus(0)) {
		return
	}

	var request struct {
		ChatID string `json:"chatId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.ChatID == "" {
		http.Error(w, "chatId is required", http.StatusBadRequest)
		return
	}

	db := mongoClient
	collection := db.Database(database.DbName).Collection(database.ChatCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"username": username,
		"chatId":   request.ChatID,
	}

	// Exclude the image field for fast loading; use projection
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: 1}}).
		SetProjection(bson.M{"image": 0})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		fmt.Println("Error loading chat history:", err)
		http.Error(w, "Failed to load chat history", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var messages []models.ChatMessage
	if err := cursor.All(ctx, &messages); err != nil {
		fmt.Println("Error decoding chat messages:", err)
		http.Error(w, "Failed to decode chat messages", http.StatusInternalServerError)
		return
	}

	if messages == nil {
		messages = []models.ChatMessage{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// LoadChatImageHandler returns the image field for a single message by its _id.
// POST /chat/image
func LoadChatImageHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	username, err := getUsernameFromToken(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !UserAuthorized(w, r, models.UserStatus(0)) {
		return
	}

	var request struct {
		MessageID string `json:"messageId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(request.MessageID)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	db := mongoClient
	collection := db.Database(database.DbName).Collection(database.ChatCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Only return image for messages owned by this user
	filter := bson.M{
		"_id":      objID,
		"username": username,
	}

	opts := options.FindOne().SetProjection(bson.M{"image": 1, "_id": 0})

	var result struct {
		Image string `bson:"image" json:"image"`
	}
	err = collection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"image": result.Image,
	})
}

// ClearChatHistoryHandler deletes all messages for a chatId + username.
// POST /chat/clear
func ClearChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	username, err := getUsernameFromToken(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !UserAuthorized(w, r, models.UserStatus(0)) {
		return
	}

	var request struct {
		ChatID string `json:"chatId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.ChatID == "" {
		http.Error(w, "chatId is required", http.StatusBadRequest)
		return
	}

	db := mongoClient
	collection := db.Database(database.DbName).Collection(database.ChatCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"username": username,
		"chatId":   request.ChatID,
	}

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		fmt.Println("Error clearing chat history:", err)
		http.Error(w, "Failed to clear chat history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deleted": result.DeletedCount,
		"success": true,
	})
}
