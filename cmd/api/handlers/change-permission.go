package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "github.com/BloxBerg-UTFPR/API-Blockchain/models"
    "github.com/BloxBerg-UTFPR/API-Blockchain/pkg/config"
    "github.com/BloxBerg-UTFPR/API-Blockchain/pkg/database"
    "go.mongodb.org/mongo-driver/bson"
    "time"
    "fmt"
    "strings"
    "strconv"
)

type ChangePermissionRequest struct {
    Username    string `json:"username"`
    Permission  string `json:"permission"`
    PermissionTime string `json:"permissionQuantity"`
    RequestAmount string `json:"requestamount"`
}

func ChangePermissionHandler(w http.ResponseWriter, r *http.Request) {
    // Check if the user is authorized (only admins can change permissions)
    if !UserAuthorized(w, r, models.UserStatus(7)) {
        return
    }

    // Parse the request body
    var request ChangePermissionRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Error decoding request body", http.StatusBadRequest)
        return
    }

    // Validate the request
    if request.Username == "" {
        http.Error(w, "Username is required", http.StatusBadRequest)
        return
    }
    if request.Permission == "" {
        http.Error(w, "Permission is required", http.StatusBadRequest)
        return
    }

    // Convert permission string to UserStatus integer
    var permissionInt int
    var permissionTimeNeeded bool
    //var permissionAmountNeeded bool
    switch strings.ToLower(request.Permission) {
    case "pending", "pendente":
        permissionInt = int(models.StatusPending)
    case "reader time based", "leitor (por tempo)":
        permissionInt = int(models.StatusReaderTimeBased)
        permissionTimeNeeded = true
    case "reader amount based", "leitor (por requisição)":
        permissionInt = int(models.StatusReaderAmountBased)
        //permissionAmountNeeded = true
    case "reader unlimited", "leitor (permanente)":
        permissionInt = int(models.StatusReaderUnlimited)
    case "editor time based", "editor (por tempo)":
        permissionInt = int(models.StatusEditorTimeBased)
        permissionTimeNeeded = true
    case "editor amount based", "editor (por requisição)":
        //permissionAmountNeeded = true
        permissionInt = int(models.StatusEditorAmountBased)
    case "editor unlimited", "editor (permanente)":
        permissionInt = int(models.StatusEditorUnlimited)
    case "admin", "administrador":
        permissionInt = int(models.StatusAdmin)
    default:
        http.Error(w, "Invalid permission level", http.StatusBadRequest)
        return
    }

    // Get the database connection
    db := database.NewMongoDB(config.MongoURI)
    collection := db.Database(database.DbName).Collection(database.CollectionName)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    ptime, err := strconv.Atoi(request.PermissionTime)

    if err != nil {
        // Handle the error (e.g., log it, return it, etc.)
        fmt.Printf("Error updating user permission: %v\n", err)
    }

    if permissionTimeNeeded && ptime <= 0 {
        http.Error(w, "Permission time is required for this permission level", http.StatusBadRequest)
        return
    }

    // if permissionAmountNeeded && ptime <= 0 {
    //     http.Error(w, "Request amount is required for this permission level", http.StatusBadRequest)
    //     return
    // }
    // Create update document
    update := bson.M{"$set": bson.M{"permission": permissionInt}}
    
    // Add time-based or amount-based parameters if provided
    if ptime > 0 {
        update["$set"].(bson.M)["accesstime"] = request.PermissionTime
    }
    
    // if strconv.Atoi(request.RequestAmount) > 0 {
    //     update["$set"].(bson.M)["reqamount"] = request.RequestAmount
    // }

    // Update the user in the database
    result, err := collection.UpdateOne(
        ctx,
        bson.M{"username": request.Username},
        update,
    )
    
    if err != nil {
        http.Error(w, "Error updating user permission", http.StatusInternalServerError)
        fmt.Printf("Error updating user permission: %v\n", err)
        return
    }

    if result.MatchedCount == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Return success response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "User permission updated successfully",
    })
}