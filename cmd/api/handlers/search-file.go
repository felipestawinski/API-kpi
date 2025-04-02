package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
	"fmt"	
    "github.com/BloxBerg-UTFPR/API-Blockchain/models"
    "github.com/BloxBerg-UTFPR/API-Blockchain/pkg/config"
    "github.com/BloxBerg-UTFPR/API-Blockchain/pkg/database"
    "go.mongodb.org/mongo-driver/bson"
)

type FileInfo struct {
	ID       int    `json:"id" bson:"id"`
	Filename string `json:"filename" bson:"filename"`
	Institution string `json:"institution" bson:"institution"`
	Writer string `json:"writer" bson:"writer"`
	Date string `json:"date" bson:"date"`
	FileAddress string `json:"fileAddress" bson:"fileAddress"`
}

func SearchFilesHandler(w http.ResponseWriter, r *http.Request) {
    //Check jwt key
    UserAuthorized(w, r)
    
    // Parse the institution from the request body
    var request struct {
        Institution string `json:"institution"`
        Id        string `json:"id"`
        FileName string `json:"filename"`
    }
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Println("request institution: ", request.Institution)
    fmt.Println("request ID: ", request.Id)
    fmt.Println("request fileName: ", request.FileName)

    if request.Institution == "" && request.Id == "" && request.FileName == "" {
        http.Error(w, "No parameter provided!", http.StatusBadRequest)
        return
    }

    // Connect to the database
    db := database.NewMongoDB(config.MongoURI)
    collection := db.Database(database.DbName).Collection(database.CollectionName)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Find all users
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        http.Error(w, "Error finding users", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var allFiles []FileInfo
    for cursor.Next(ctx) {
        var user models.User
        if err := cursor.Decode(&user); err != nil {
            fmt.Printf("Error decoding user: %v\n", err)
            continue
        }

        // Parse files for this user
        for _, fileStr := range user.Files {
            var fileInfo FileInfo
            if err := json.Unmarshal([]byte(fileStr), &fileInfo); err != nil {
                fmt.Printf("Error parsing file info: %v\n", err)
                continue
            }
            
            // Only include files that match the requested institution (if given)
            if fileInfo.Institution == request.Institution {
                allFiles = append(allFiles, fileInfo)
            }

            // Only include files that match the requested filename (if given)
            if fileInfo.Filename == request.FileName {
                allFiles = append(allFiles, fileInfo)
            }
        }
    }

    if len(allFiles) == 0 {
        http.Error(w, "No files found for this institution", http.StatusNotFound)
        return
    }

    // Return all matching files in JSON format
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string][]FileInfo{
        "files": allFiles,
    })
}