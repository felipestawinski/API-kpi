package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/felipestawinski/API-kpi/cmd/api/handlers"
	"github.com/felipestawinski/API-kpi/pkg/database"
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
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file couldn't be loaded")
	}

	// Initialize the shared MongoDB client once at startup.
	mongoClient := database.NewMongoDB(database.MongoURI)
	database.SetClient(mongoClient)
	handlers.SetMongoClient(mongoClient)

	// Gracefully disconnect the client when main returns.
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Println("MongoDB disconnect error:", err)
		}
	}()

	// Ensure indexes now that the shared client is ready.
	database.EnsureChatIndexes()

	fmt.Println("Server starting on port 8080")

	// Rotas
	http.Handle("/register", enableCORS(http.HandlerFunc(handlers.RegisterHandler)))
	http.Handle("/login", enableCORS(http.HandlerFunc(handlers.LoginHandler)))
	http.Handle("/upload", enableCORS(http.HandlerFunc(handlers.UploadFileHandler)))
	http.Handle("/health-check", enableCORS(http.HandlerFunc(handlers.HealthCheckOnlyHandler)))
	http.Handle("/upload-confirmed", enableCORS(http.HandlerFunc(handlers.UploadConfirmedHandler)))
	http.Handle("/upload-picture", enableCORS(http.HandlerFunc(handlers.UploadPictureHandler)))
	http.Handle("/files", enableCORS(http.HandlerFunc(handlers.GetFilesHandler)))
	http.Handle("/search-file", enableCORS(http.HandlerFunc(handlers.SearchFilesHandler)))
	http.Handle("/user-info", enableCORS(http.HandlerFunc(handlers.UserInfoHandler)))
	http.Handle("/users", enableCORS(http.HandlerFunc(handlers.GetUsersHandler)))
	http.Handle("/pending-users", enableCORS(http.HandlerFunc(handlers.GetPendingUsersHandler)))
	http.Handle("/change-permission", enableCORS(http.HandlerFunc(handlers.ChangePermissionHandler)))
	http.Handle("/analysis-gen", enableCORS(http.HandlerFunc(handlers.AnalysisGenHandler)))
	http.Handle("/file-preview", enableCORS(http.HandlerFunc(handlers.FilePreviewHandler)))
	http.Handle("/token-usage", enableCORS(http.HandlerFunc(handlers.TokenUsageHandler)))

	// Chat message persistence routes
	http.Handle("/chat/save", enableCORS(http.HandlerFunc(handlers.SaveChatMessageHandler)))
	http.Handle("/chat/load", enableCORS(http.HandlerFunc(handlers.LoadChatHistoryHandler)))
	http.Handle("/chat/image", enableCORS(http.HandlerFunc(handlers.LoadChatImageHandler)))
	http.Handle("/chat/clear", enableCORS(http.HandlerFunc(handlers.ClearChatHistoryHandler)))

	// Gallery routes
	http.Handle("/gallery/save", enableCORS(http.HandlerFunc(handlers.SaveToGalleryHandler)))
	http.Handle("/gallery/load", enableCORS(http.HandlerFunc(handlers.LoadGalleryHandler)))
	http.Handle("/gallery/delete", enableCORS(http.HandlerFunc(handlers.DeleteGalleryImageHandler)))

	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
