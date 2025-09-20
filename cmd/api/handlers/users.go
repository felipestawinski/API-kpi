package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "github.com/felipestawinski/API-kpi/models"
    "github.com/felipestawinski/API-kpi/pkg/config"
    "github.com/felipestawinski/API-kpi/pkg/database"
    "go.mongodb.org/mongo-driver/bson"
    "time"
    "fmt"
)

type UserResponse struct {
    Email       string   `json:"email"`
    Password    string   `json:"password"`
    Username    string   `json:"username"`
    Institution string   `json:"institution"`
    Role        string   `json:"role"`
    Permission  string   `json:"permission"`
    ID          string   `json:"id,omitempty"`
    Files       []string `json:"files,omitempty"`
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {


    // Check if the user is authorized
    if !UserAuthorized(w, r, models.UserStatus(0)) {
        return 
    }

    // Get the list of users from the database
    db := database.NewMongoDB(config.MongoURI)
    collection := db.Database(database.DbName).Collection(database.CollectionName)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Get username from JWT token
	tokenStr := r.Header.Get("Authorization")
	username, err := getUsernameFromToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

    fmt.Println("Requesting user list for:", username)

    cursor, err := collection.Find(ctx, bson.M{})
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
        
        // Create response with string permission
        userResponse := UserResponse{
            ID:          user.ID,
            Email:       user.Email,
            Username:    user.Username,
            Institution: user.Institution,

        }
        userResponses = append(userResponses, userResponse)
    }

    // Encode the users as JSON and send the response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userResponses)
}