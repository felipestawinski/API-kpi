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
)

type UserPendingResponse struct {
    Email       string   `json:"email"`
    Password    string   `json:"password"`
    Username    string   `json:"username"`
    Institution string   `json:"institution"`
    Role        string   `json:"role"`
    Permission  string   `json:"permission"`
    ID          string   `json:"id,omitempty"`
}

// UserRequest represents the request body for the GetUsersHandler
type UserPendingRequest struct {
    Username string `json:"username"`
}

func GetPendingUsersHandler(w http.ResponseWriter, r *http.Request) {

    // Check if the user is authorized
    UserAuthorized(w, r, models.UserStatus(6))

    // Parse the request body
    var request UserRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Error decoding request body", http.StatusBadRequest)
        return
    }

    // Get the list of users from the database
    db := database.NewMongoDB(config.MongoURI)
    collection := db.Database(database.DbName).Collection(database.CollectionName)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Only return pending users
	filter := bson.M{"permission": 0}
	fmt.Println("Filtering for pending users only")

    // Find returns a cursor and an error
    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        http.Error(w, "Error retrieving users", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    // Decode all documents from the cursor into a slice of users
    var users []models.User
    if err := cursor.All(ctx, &users); err != nil {
        http.Error(w, "Error decoding users", http.StatusInternalServerError)
        return
    }

    // Convert User to UserResponse with string permissions
    var userResponses []UserResponse
    for _, user := range users {
        // Convert int permission to UserStatus type
        permissionStatus := models.UserStatus(user.Permission)
        
        // Create response with string permission
        userResponse := UserResponse{
            Email:       user.Email,
            Password:    user.Password,
            Username:    user.Username,
            Institution: user.Institution,
            Role:        user.Role,
            Permission:  permissionStatus.String(),
            ID:          user.ID,
        }
        userResponses = append(userResponses, userResponse)
    }

    // Encode the users as JSON and send the response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userResponses)
}