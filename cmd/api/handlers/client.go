package handlers

import "go.mongodb.org/mongo-driver/mongo"

// mongoClient is the shared *mongo.Client injected at application startup.
// All handler functions use this instead of creating a new connection per request.
var mongoClient *mongo.Client

// SetMongoClient registers the application-wide MongoDB client.
// Must be called from main() before any HTTP requests are served.
func SetMongoClient(c *mongo.Client) {
	mongoClient = c
}
