package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"mime/multipart"
	"os"
	"github.com/BloxBerg-UTFPR/API-Blockchain/pkg/config"
	"github.com/BloxBerg-UTFPR/API-Blockchain/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"github.com/BloxBerg-UTFPR/API-Blockchain/models"
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

// uploadFileHandler handles the HTTP request for file upload
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {


	//UserAuthorized(w, r)
	
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	fmt.Println("username", username)
    if username == "" {
        http.Error(w, "Username is required", http.StatusBadRequest)
        return
    }

	institution := r.FormValue("institution")
	fmt.Println("institution", institution)
    if institution == "" {
        http.Error(w, "Institution is required", http.StatusBadRequest)
        return
    }

	db := database.NewMongoDB(config.MongoURI)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	collection := db.Database(database.DbName).Collection(database.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	fmt.Printf("Request received: %+v\n", r)

	fmt.Printf("Request received:\nMethod: %s\nHeaders: %v\nContent-Type: %s\n", 
		r.Method, r.Header, r.Header.Get("Content-Type"))


	err = r.ParseMultipartForm(10 << 20) // Limit your max input length!
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

	uri := "https://scarlet-implicit-lobster-990.mypinata.cloud/ipfs/" + ipfsHash


	// Format the blockchain POST URL
	blockchainURL := fmt.Sprintf("http://localhost:8080/blockchain/PostData?entity=%s&uri=%s", 
	url.QueryEscape(institution), 
	url.QueryEscape(uri))

	// Make the POST request
	resp, err := http.Post(blockchainURL, "application/json", nil)
	if err != nil {
	http.Error(w, fmt.Sprintf("Error calling blockchain: %v", err), http.StatusInternalServerError)
	return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
	http.Error(w, fmt.Sprintf("Error reading response: %v", err), http.StatusInternalServerError)
	return
	}
	fmt.Printf("body: %s", body)

	// Create file info struct
	type FileInfo struct {
	ID       int    `json:"id" bson:"id"`
	Filename string `json:"filename" bson:"filename"`
	Institution string `json:"institution" bson:"institution"`
	ContractAddress string `json:"contractAddress" bson:"contractAddress"`
	TxHash   string `json:"txHash" bson:"txHash"`
	IfpsHash string `json:"ifpsHash" bson:"ifpsHash"`
	}

	// Define struct for JSON response
	type BlockchainResponse struct {
		TxHash string `json:"txHash"`
	}

	// Parse JSON response
	var response BlockchainResponse
	if err := json.Unmarshal(body, &response); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing response: %v", err), http.StatusInternalServerError)
		return
	}

	// Get clean hash value
	txHash := response.TxHash
	fmt.Printf("txHash: %s\n", txHash)
	contractAddress := "0x473f8eA5Ce1F35acf7Eb61A6D4b74C8f5cf2f362"

	// Determine new ID
	newID := 1
	if len(user.Files) > 0 {
		// Find highest ID
		maxID := 0
		for _, file := range user.Files {
			var fileInfo FileInfo
			if err := json.Unmarshal([]byte(file), &fileInfo); err != nil {
				continue
			}
			if fileInfo.ID > maxID {
				maxID = fileInfo.ID
			}
		}
		newID = maxID + 1
	}

	// Create new file entry
	newFile := FileInfo{
		ID:       newID,
		Filename: handler.Filename,
		Institution: institution,
		ContractAddress: contractAddress,
		TxHash:   txHash,
		IfpsHash: ipfsHash,
	}

	fmt.Printf("Inserting new file: %+v\n", newFile)

	// Convert to JSON string
	newFileJSON, err := json.Marshal(newFile)
	if err != nil {
		http.Error(w, "Error creating file entry", http.StatusInternalServerError)
		return
	}

	// Update user document
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"username": username},
		bson.M{
			"$push": bson.M{
				"files": string(newFileJSON),
			},
		},
	)
	if err != nil {
		http.Error(w, "Error updating user files", http.StatusInternalServerError)
		return
	}


	// Upload the JSON file to IPFS
	// jsonFile, err := os.Open(jsonFilePath)
	// if err != nil {
	// 	http.Error(w, "Error opening JSON file", http.StatusInternalServerError)
	// 	return
	// }
	// defer jsonFile.Close()
	// //jsonIpfsHash, err := uploadFileToPinata(jsonFile, handler.Filename+".json")
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Error uploading JSON file: %v", err), http.StatusInternalServerError)
	// 	return
	// }

	// Respond with both IPFS hashes
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"fileId": fmt.Sprintf("%d", newID), //convert to int
	})
}