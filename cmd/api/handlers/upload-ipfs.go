package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/felipestawinski/API-kpi/models"
	"github.com/felipestawinski/API-kpi/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var allowedFileTypes = map[string]bool{
	"csv":  true,
	"xlsx": true,
	"json": true,
}

func detectFileType(filename string, contentType string) string {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")
	if allowedFileTypes[ext] {
		return ext
	}

	if contentType == "" {
		return ""
	}

	parsedType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		parsedType = strings.ToLower(contentType)
	}

	switch {
	case strings.Contains(parsedType, "csv") || strings.Contains(parsedType, "text/plain"):
		return "csv"
	case strings.Contains(parsedType, "spreadsheet") || strings.Contains(parsedType, "excel"):
		return "xlsx"
	case strings.Contains(parsedType, "json"):
		return "json"
	default:
		return ""
	}
}

func ensureFilenameHasExtension(baseName string, sourceFilename string, fileType string) string {
	if strings.TrimSpace(baseName) == "" {
		baseName = sourceFilename
	}

	if filepath.Ext(baseName) != "" {
		return baseName
	}

	if fileType == "" {
		return baseName
	}

	return fmt.Sprintf("%s.%s", baseName, fileType)
}

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

	db := mongoClient
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
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file from request", http.StatusBadRequest)
		fmt.Printf("Unable to read file from request: %v\n", err)
		return
	}
	defer file.Close()

	fileType := detectFileType(fileHeader.Filename, fileHeader.Header.Get("Content-Type"))
	if fileType == "" {
		http.Error(w, "Unsupported file type. Allowed types: csv, xlsx, json", http.StatusBadRequest)
		return
	}

	filename = ensureFilenameHasExtension(filename, fileHeader.Filename, fileType)

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

	// Pre-warm the Python DataFrame cache in the background (fire-and-forget)
	fireAndForgetPreload(uri, fileType)

	// Create file info struct
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
	newFile := models.File{
		ID:          newID,
		Filename:    filename,
		Institution: r.FormValue("institution"),
		Writer:      username,
		Date:        time.Now().Format("2006-01-02"),
		FileAddress: uri,
		FileType:    fileType,
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

// HealthCheckOnlyHandler runs the data health check WITHOUT uploading to IPFS.
// The frontend calls this first so the user can review the report before deciding.
func HealthCheckOnlyHandler(w http.ResponseWriter, r *http.Request) {
	if !UserAuthorized(w, r, models.UserStatus(0)) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusInternalServerError)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileType := detectFileType(fileHeader.Filename, fileHeader.Header.Get("Content-Type"))
	if fileType == "" {
		http.Error(w, "Unsupported file type. Allowed types: csv, xlsx, json", http.StatusBadRequest)
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file data", http.StatusInternalServerError)
		return
	}

	filename := r.FormValue("filename")
	if filename == "" {
		filename = ensureFilenameHasExtension("file", fileHeader.Filename, fileType)
	} else {
		filename = ensureFilenameHasExtension(filename, fileHeader.Filename, fileType)
	}

	analysis, err := sendToDataHealthCheck(fileData, filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Health check failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"dataHealthCheck": analysis,
	})
}

// sendToDataHealthCheckClean sends a file to the Python /data-health-check-clean
// endpoint and returns the cleaned CSV bytes.
func sendToDataHealthCheckClean(fileData []byte, filename string) ([]byte, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file for clean: %v", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(fileData)); err != nil {
		return nil, fmt.Errorf("failed to copy file for clean: %v", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer for clean: %v", err)
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:9090/data-health-check-clean", &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create clean request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send file to clean service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("clean service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	cleanedData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read cleaned file response: %v", err)
	}

	return cleanedData, nil
}

// fireAndForgetPreload asynchronously notifies the Python analysis service
// to pre-parse and cache the uploaded file as a DataFrame.
// Errors are logged but never surface to the caller.
func fireAndForgetPreload(fileAddress string, fileType string) {
	go func() {
		payload, err := json.Marshal(map[string]string{
			"fileAddress": fileAddress,
			"fileType":    fileType,
		})
		if err != nil {
			fmt.Printf("[preload] Failed to marshal payload: %v\n", err)
			return
		}

		client := &http.Client{Timeout: 120 * time.Second}
		resp, err := client.Post(
			"http://127.0.0.1:9090/preload-file",
			"application/json",
			bytes.NewReader(payload),
		)
		if err != nil {
			fmt.Printf("[preload] Request failed for %s: %v\n", fileAddress, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("[preload] Non-OK response for %s: status=%d body=%s\n", fileAddress, resp.StatusCode, string(body))
			return
		}
		fmt.Printf("[preload] Successfully preloaded: %s\n", fileAddress)
	}()
}

// UploadConfirmedHandler uploads the file to IPFS after the user has reviewed
// the health check. Accepts a "mode" form field: "raw" or "cleaned".
func UploadConfirmedHandler(w http.ResponseWriter, r *http.Request) {
	if !UserAuthorized(w, r, models.UserStatus(0)) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	filename := r.FormValue("filename")
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

	db := mongoClient
	collection := db.Database(database.DbName).Collection(database.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.User
	err = collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusInternalServerError)
		return
	}

	mode := r.FormValue("mode") // "raw" or "cleaned"
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileType := detectFileType(fileHeader.Filename, fileHeader.Header.Get("Content-Type"))
	if fileType == "" {
		http.Error(w, "Unsupported file type. Allowed types: csv, xlsx, json", http.StatusBadRequest)
		return
	}

	filename = ensureFilenameHasExtension(filename, fileHeader.Filename, fileType)

	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file data", http.StatusInternalServerError)
		return
	}

	// Determine which data to upload
	var uploadData []byte
	if mode == "cleaned" {
		fmt.Println("User chose cleaned file, requesting cleaning from analysis-gen...")
		cleanedData, err := sendToDataHealthCheckClean(fileData, filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error cleaning file: %v", err), http.StatusInternalServerError)
			return
		}
		uploadData = cleanedData
	} else {
		// Default to raw
		uploadData = fileData
	}

	// Upload to IPFS
	ipfsHash, err := uploadFileToPinata(bytes.NewReader(uploadData), filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading file: %v", err), http.StatusInternalServerError)
		return
	}

	uri := "https://scarlet-implicit-lobster-990.mypinata.cloud/ipfs/" + ipfsHash

	// Pre-warm the Python DataFrame cache in the background (fire-and-forget)
	fireAndForgetPreload(uri, fileType)

	newID := 1
	if len(user.Files) > 0 {
		maxID := 0
		for _, file := range user.Files {
			if file.ID > maxID {
				maxID = file.ID
			}
		}
		newID = maxID + 1
	}

	newFile := models.File{
		ID:          newID,
		Filename:    filename,
		Institution: r.FormValue("institution"),
		Writer:      username,
		Date:        time.Now().Format("2006-01-02"),
		FileAddress: uri,
		FileType:    fileType,
	}

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"fileId": fmt.Sprintf("%d", newID),
	})
}
