package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"mime/multipart"
	"os"
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

// uploadFileHandler handles the HTTP request for file upload
// uploadFileHandler handles the HTTP request for file upload
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	UserAuthorized(w, r)
	
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Request received")
	err := r.ParseMultipartForm(10 << 20) // Limit your max input length!
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusInternalServerError)
		fmt.Printf("Error parsing multipart form: %v\n", err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file from request", http.StatusBadRequest)
		fmt.Printf("Unable to read file from request: %v\n", err)
		return
	}
	defer file.Close()

	// Create JSON file with file details
	fileDetails := map[string]interface{}{
		"name": handler.Filename,
		"size": handler.Size,
	}
	fileDetailsBytes, err := json.Marshal(fileDetails)
	if err != nil {
		http.Error(w, "Error creating file details JSON", http.StatusInternalServerError)
		return
	}
	jsonFilePath := "json-files/" + handler.Filename + ".json"
	err = os.WriteFile(jsonFilePath, fileDetailsBytes, 0644)
	if err != nil {
		http.Error(w, "Error saving file details JSON", http.StatusInternalServerError)
		return
	}

	// Upload the original file to IPFS
	ipfsHash, err := uploadFileToPinata(file, handler.Filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading file: %v", err), http.StatusInternalServerError)
		fmt.Printf("Error uploading file to Pinata: %v\n", err)
		return
	}

	// Upload the JSON file to IPFS
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		http.Error(w, "Error opening JSON file", http.StatusInternalServerError)
		return
	}
	defer jsonFile.Close()
	jsonIpfsHash, err := uploadFileToPinata(jsonFile, handler.Filename+".json")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading JSON file: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with both IPFS hashes
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"fileIpfsHash": ipfsHash,
		"jsonIpfsHash": jsonIpfsHash,
	})
}