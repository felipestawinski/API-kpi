package main
import (
	"fmt"
	"log"
	"net/http"
	"github.com/BloxBerg-UTFPR/API-Blockchain/cmd/api/handlers"
	"github.com/joho/godotenv"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file couldn't be loaded")
	}
	// Rotas
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/upload", handlers.UploadFileHandler)
	http.HandleFunc("/blockchain/{method}", handlers.BlockchainInteraction)

	log.Fatal(http.ListenAndServe(":8080", nil))

	fmt.Println("Server starting on port 8080")
}

