package main
import (
	//"context"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/BloxBerg-UTFPR/API-Blockchain/cmd/api/handlers"
	"github.com/BloxBerg-UTFPR/API-Blockchain/cmd/api/blockchain"
)


func main() {

	// Rotas
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/welcome", handlers.WelcomeHandler)
	http.HandleFunc("/upload", handlers.UploadFileHandler)
	blockchain.Blockchain()

	log.Fatal(http.ListenAndServe(":8080", nil))

	fmt.Println("Server starting on port 8080")
}

