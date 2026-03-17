package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"mime/multipart"
	"os"
	"strings"
	"github.com/felipestawinski/API-kpi/pkg/config"
	"github.com/felipestawinski/API-kpi/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"github.com/felipestawinski/API-kpi/models"
	"time" 
)

// uploadFileToPinata uploads a file to Pinata and returns the IPFS hash
func uploadFileToPinata(file io.Reader, filename string) (string, error) {

	
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

	fmt.Println("apiKey->", apiKey)
	fmt.Println("apiSecret->", apiSecret)
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

// sendToDataHealthCheck sends a file to the data health check service and returns the analysis
func sendToDataHealthCheck(fileData []byte, filename string) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file for health check: %v", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(fileData)); err != nil {
		return "", fmt.Errorf("failed to copy file for health check: %v", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer for health check: %v", err)
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:9090/data-health-check", &body)
	if err != nil {
		return "", fmt.Errorf("failed to create health check request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send file to health check service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("health check service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read health check response: %v", err)
	}

	return string(respBody), nil
}

// uploadFileHandler handles the HTTP request for file upload
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {

	//Check jwt key
	if !UserAuthorized(w, r, models.UserStatus(0)) {
		return
	}
	
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	filename := r.FormValue("filename")
	fmt.Println("filename: ", filename)
    if filename == "" {
        http.Error(w, "Filename is required", http.StatusBadRequest)
        return
    }

	tokenStr := r.Header.Get("Authorization")
	username, err := getUsernameFromToken(tokenStr)

	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	db := database.NewMongoDB(config.MongoURI)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	collection := db.Database(database.DbName).Collection(database.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.User
	err = collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	fmt.Printf("Uploading file to IPFS")

	err = r.ParseMultipartForm(10 << 20) // Limit your max input length
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusInternalServerError)
		fmt.Printf("Error parsing multipart form: %v\n", err)
		return
	}

	// Read dataHealthCheck flag from form
	dataHealthCheck := strings.EqualFold(r.FormValue("dataHealthCheck"), "true")
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file from request", http.StatusBadRequest)
		fmt.Printf("Unable to read file from request: %v\n", err)
		return
	}
	defer file.Close()


	// Buffer the file so it can be read multiple times (IPFS upload + health check)
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file data", http.StatusInternalServerError)
		fmt.Printf("Error reading file data: %v\n", err)
		return
	}


	// Upload the original file to IPFS (using buffered data)
	ipfsHash, err := uploadFileToPinata(bytes.NewReader(fileData), filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading file: %v", err), http.StatusInternalServerError)
		fmt.Printf("Error uploading file to Pinata: %v\n", err)
		return
	}

	// If dataHealthCheck is enabled, send the file for analysis
	var healthCheckAnalysis string
	if dataHealthCheck {
		fmt.Println("Data health check enabled, sending file for analysis...")
		analysis, err := sendToDataHealthCheck(fileData, filename)
		if err != nil {
			fmt.Printf("Warning: data health check failed: %v\n", err)
			healthCheckAnalysis = fmt.Sprintf("Health check failed: %v", err)
		} else {
			healthCheckAnalysis = analysis
		}
	}

	uri := "https://scarlet-implicit-lobster-990.mypinata.cloud/ipfs/" + ipfsHash



	// Create file info struct
	type FileInfo struct {
		ID       int    `json:"id" bson:"id"`
		Filename string `json:"filename" bson:"filename"`
		Institution string `json:"institution" bson:"institution"`
		Writer string `json:"writer" bson:"writer"`
		Date string `json:"date" bson:"date"`
		FileAddress string `json:"fileAddress" bson:"fileAddress"`
	}
	

	// Determine new ID
	newID := 1
	if len(user.Files) > 0 {
		// Find highest ID
		maxID := 0
		for _, file := range user.Files {
			
			if file.ID > maxID {
				maxID = file.ID
			}
		}
		newID = maxID + 1
	}

	// Create new file entry
	newFile := FileInfo{
		ID:       newID,
		Filename: filename,
		Institution: r.FormValue("institution"),
		Writer: username,
		Date: time.Now().Format("2006-01-02"),
		FileAddress: uri,
	}

	fmt.Printf("Inserting new file: %+v\n", newFile)


	// Update user document
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"username": username},
		bson.M{
			"$push": bson.M{
				"files": newFile,
			},
		},
	)
	if err != nil {
		http.Error(w, "Error updating user files", http.StatusInternalServerError)
		return
	}


	// Respond with the file ID and optional health check analysis
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"fileId": fmt.Sprintf("%d", newID),
	}
	if dataHealthCheck {
		response["dataHealthCheck"] = healthCheckAnalysis
	}
	json.NewEncoder(w).Encode(response)
}
