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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveToGalleryHandler saves a base64 image to the user's gallery.
// POST /gallery/save
func SaveToGalleryHandler(w http.ResponseWriter, r *http.Request) {
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
		Image string `json:"image"`
		Title string `json:"title,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.Image == "" {
		http.Error(w, "image is required", http.StatusBadRequest)
		return
	}

	galleryImage := models.GalleryImage{
		Username:  username,
		Image:     request.Image,
		Title:     request.Title,
		CreatedAt: time.Now(),
	}

	db := database.NewMongoDB(config.MongoURI)
	collection := db.Database(database.DbName).Collection(database.GalleryCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, galleryImage)
	if err != nil {
		fmt.Println("Error saving gallery image:", err)
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"_id":     result.InsertedID,
		"success": true,
	})
}

// LoadGalleryHandler returns all gallery images for the authenticated user.
// POST /gallery/load
func LoadGalleryHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	username, err := getUsernameFromToken(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !UserAuthorized(w, r, models.UserStatus(0)) {
		return
	}

	db := database.NewMongoDB(config.MongoURI)
	collection := db.Database(database.DbName).Collection(database.GalleryCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		fmt.Println("Error loading gallery:", err)
		http.Error(w, "Failed to load gallery", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var images []models.GalleryImage
	if err := cursor.All(ctx, &images); err != nil {
		fmt.Println("Error decoding gallery images:", err)
		http.Error(w, "Failed to decode gallery images", http.StatusInternalServerError)
		return
	}

	if images == nil {
		images = []models.GalleryImage{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

// DeleteGalleryImageHandler deletes a single gallery image by its _id.
// POST /gallery/delete
func DeleteGalleryImageHandler(w http.ResponseWriter, r *http.Request) {
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
		ImageID string `json:"imageId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(request.ImageID)
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	db := database.NewMongoDB(config.MongoURI)
	collection := db.Database(database.DbName).Collection(database.GalleryCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":      objID,
		"username": username,
	}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("Error deleting gallery image:", err)
		http.Error(w, "Failed to delete image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deleted": result.DeletedCount,
		"success": true,
	})
}
