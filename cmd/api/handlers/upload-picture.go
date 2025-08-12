package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
	"github.com/felipestawinski/API-kpi/models"
	"github.com/felipestawinski/API-kpi/pkg/config"
	"github.com/felipestawinski/API-kpi/pkg/database"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

// uploadPictureToPinata uploads a picture file to Pinata and returns the IPFS hash
func uploadPictureToPinata(file io.Reader, filename string) (string, error) {
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	// Prepare the form data
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("failed to copy file to form: %v", err)
	}

	// Close the writer to finalize the form data
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinFileToIPFS", &body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("pinata_api_key", apiKey)
	req.Header.Set("pinata_secret_api_key", apiSecret)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad status code: %v, response: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", fmt.Errorf("unable to parse response: %v", err)
	}

	ipfsHash, ok := result["IpfsHash"].(string)
	if !ok {
		return "", fmt.Errorf("unable to find IpfsHash in response")
	}

	return ipfsHash, nil
}

// getUsernameFromToken extracts username from JWT token
func getUsernameFromToken(tokenStr string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	return claims.Subject, nil
}

// UploadPictureHandler handles the HTTP request for profile picture upload
func UploadPictureHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authorized (minimum permission level 1 - any authenticated user)
	if !UserAuthorized(w, r, models.UserStatus(0)) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get username from JWT token
	tokenStr := r.Header.Get("Authorization")
	username, err := getUsernameFromToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(10 << 20) // Limit to 10MB
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		fmt.Printf("Error parsing multipart form: %v\n", err)
		return
	}

	// Debug: Print available form fields
	fmt.Printf("Available form fields: %v\n", r.MultipartForm.File)
	fmt.Printf("Form values: %v\n", r.MultipartForm.Value)

	file, handler, err := r.FormFile("profilePicture")
	if err != nil {
		http.Error(w, "Unable to read picture from request", http.StatusBadRequest)
		fmt.Printf("Unable to read picture from request: %v\n", err)
		return
	}
	defer file.Close()

	// Validate file type (basic check)
	contentType := handler.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" && contentType != "image/gif" {
		http.Error(w, "Invalid file type. Only JPEG, PNG, and GIF images are allowed", http.StatusBadRequest)
		return
	}

	// Validate file size (max 5MB)
	if handler.Size > 5*1024*1024 {
		http.Error(w, "File too large. Maximum size is 5MB", http.StatusBadRequest)
		return
	}

	fmt.Printf("Uploading profile picture for user: %s\n", username)

	// Upload the picture to IPFS
	ipfsHash, err := uploadPictureToPinata(file, fmt.Sprintf("%s_profile_%s", username, handler.Filename))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading picture: %v", err), http.StatusInternalServerError)
		fmt.Printf("Error uploading picture to Pinata: %v\n", err)
		return
	}

	// Create the full URI for the picture
	pictureURI := "https://scarlet-implicit-lobster-990.mypinata.cloud/ipfs/" + ipfsHash

	// Update user's profile picture in database
	db := database.NewMongoDB(config.MongoURI)
	collection := db.Database(database.DbName).Collection(database.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Update user document with new profile picture URI
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"username": username},
		bson.M{
			"$set": bson.M{
				"profilePicture": pictureURI,
			},
		},
	)
	if err != nil {
		http.Error(w, "Error updating user profile picture", http.StatusInternalServerError)
		fmt.Printf("Error updating user profile picture: %v\n", err)
		return
	}

	fmt.Printf("Profile picture updated successfully for user: %s\n", username)

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":        true,
		"message":        "Profile picture uploaded successfully",
		"profilePicture": pictureURI,
		"ipfsHash":       ipfsHash,
	})
}
