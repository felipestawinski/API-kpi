package main
import (
	//"context"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/BloxBerg-UTFPR/API-Blockchain/cmd/api/handlers"
)
// JWT secret key *TODO: Armazenar essa chave em um lugar seguro
// var jwtKey = []byte("my_secret_key")

// const mongoURI = "mongodb://localhost:27017"
// const dbName = "userdb"
// const collectionName = "users"
// var client *mongo.Client


func main() {
	// rotas
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/welcome", handlers.WelcomeHandler)
	//http.HandleFunc("/logout", logoutHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

	fmt.Println("Server starting on port 8080")
}

