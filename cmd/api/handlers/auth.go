package handlers

import (
	"fmt"
	"net/http"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"github.com/BloxBerg-UTFPR/API-Blockchain/pkg/database"
	"github.com/BloxBerg-UTFPR/API-Blockchain/pkg/config"
	"github.com/BloxBerg-UTFPR/API-Blockchain/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)


func UserAuthorized(w http.ResponseWriter, r *http.Request, permissionLevel models.UserStatus) bool {

	//Check for the Authorization header
	tokenStr := r.Header.Get("Authorization")
	println("Token jwt: ", tokenStr)
	if tokenStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	client := database.NewMongoDB(config.MongoURI)

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		fmt.Println("Erro: chave jwt inválida", token)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	collection := client.Database(database.DbName).Collection(database.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find user by username
	var user models.User
	err = collection.FindOne(ctx, bson.M{"username": claims.Subject}).Decode(&user)
	if err != nil {
		fmt.Println("Erro: usuário não encontrado", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	println("permission required to execute action:", permissionLevel)
	println("user permission:", user.Permission)

	if int(permissionLevel) > user.Permission {
		fmt.Println("Erro: permissão insuficiente")
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return false
	}

	// Check permission level
	// switch user.Permission {
	// 	case StatusPending:
	// 		http.Error(w, "Pending", http.StatusForbidden)
	// 		return
	// 	case StatusReaderTimeBased:
	// 		if permissionLevel 
	// }


	return true
}