package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time" 
	"go.mongodb.org/mongo-driver/bson"
	"github.com/BloxBerg-UTFPR/API-Blockchain/models"
	"github.com/BloxBerg-UTFPR/API-Blockchain/pkg/config"
	"github.com/BloxBerg-UTFPR/API-Blockchain/pkg/database"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"os"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

//TODO* Armazenar a key em um lugar seguro
var jwtKey = []byte(os.Getenv("JWT_KEY"))
// Função chamada para a rota /register

// CORS middleware
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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	db := database.NewMongoDB(config.MongoURI)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("user", user)

	// Validate password
	if err := validatePassword(user.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := db.Database(database.DbName).Collection(database.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verificação - usuário já existe
	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Verificação - email ja cadastrado
	var userEmail models.User
	err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&userEmail)
	if err == nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// Criptografia da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Default permission level
	user.Permission = 3

	// Insere novo usuário
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// After successful user registration
	if err := sendWelcomeEmail(user.Email); err != nil {
		// Log the error, but don't return it to the user
		fmt.Printf("Error sending welcome email: %v\n", err)
	}


	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully")
}

// validatePassword checks if the password meets the required criteria
func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}
	return nil
}

func sendWelcomeEmail(email string) error {
	fmt.Printf("Sending email to %s", email)
	from := mail.NewEmail("Blockchain", "felipe.stawinski@gmail.com")
	subject := "Welcome"
	to := mail.NewEmail("New User", email)
	plainTextContent := "Welcome! We're glad to have you on board."
	htmlContent := "<strong>Welcome to Your App!</strong><p>We're glad to have you on board.</p>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	return err
}