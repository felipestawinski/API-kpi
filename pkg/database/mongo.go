package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MongoURI = "mongodb://localhost:27017"
const DbName = "kpidb"
const CollectionName = "users"
const ChatCollectionName = "chat_messages"
const GalleryCollectionName = "gallery_images"

var client *mongo.Client

// NewMongoDB initializes a new MongoDB client.
func NewMongoDB(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

// EnsureChatIndexes creates the compound index on chat_messages for fast lookups.
func EnsureChatIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := NewMongoDB(MongoURI)
	collection := db.Database(DbName).Collection(ChatCollectionName)

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "username", Value: 1},
			{Key: "chatId", Value: 1},
			{Key: "timestamp", Value: 1},
		},
		Options: options.Index().SetName("idx_username_chatId_timestamp"),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		fmt.Println("Warning: failed to create chat index:", err)
	} else {
		fmt.Println("Chat indexes ensured successfully")
	}
}
