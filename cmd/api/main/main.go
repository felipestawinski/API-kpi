package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/felipestawinski/API-kpi/cmd/api/handlers"
	"github.com/joho/godotenv"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
	http.Handle("/login", enableCORS(http.HandlerFunc(handlers.LoginHandler)))
	http.Handle("/upload", enableCORS(http.HandlerFunc(handlers.UploadFileHandler)))
	http.Handle("/files", enableCORS(http.HandlerFunc(handlers.GetFilesHandler)))
	http.Handle("/search-file", enableCORS(http.HandlerFunc(handlers.SearchFilesHandler)))
	http.Handle("/user-info", enableCORS(http.HandlerFunc(handlers.UserInfoHandler)))
	http.Handle("/blockchain/{method}", enableCORS(http.HandlerFunc(handlers.BlockchainInteraction)))
	http.Handle("/users", enableCORS(http.HandlerFunc(handlers.GetUsersHandler)))
	http.Handle("/pending-users", enableCORS(http.HandlerFunc(handlers.GetPendingUsersHandler)))
	http.Handle("/change-permission", enableCORS(http.HandlerFunc(handlers.ChangePermissionHandler)))

	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}