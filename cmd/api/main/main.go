package main
import (
	"fmt"
	"log"
	"net/http"
	"github.com/BloxBerg-UTFPR/API-Blockchain/cmd/api/handlers"
	"github.com/joho/godotenv"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
    // Add error handling for server startup
    err := godotenv.Load()
    if err != nil {
        log.Fatal(".env file couldn't be loaded")
    }

    // Move print before ListenAndServe
    fmt.Println("Server starting on port 8080")

    // Rotas
    http.Handle("/register", enableCORS(http.HandlerFunc(handlers.RegisterHandler)))
    http.HandleFunc("/login", handlers.LoginHandler)
    http.HandleFunc("/upload", handlers.UploadFileHandler)
    http.HandleFunc("/files", handlers.GetFilesHandler)
	http.HandleFunc("/search-file", handlers.SearchFilesHandler)
	http.HandleFunc("/user-info", handlers.UserInfoHandler)
    http.HandleFunc("/blockchain/{method}", handlers.BlockchainInteraction)

    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}

