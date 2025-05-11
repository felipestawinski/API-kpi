package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MongoURI = "mongodb+srv://edfisic:DwyZQ0wzP3veQEWf@db-1.5di98i7.mongodb.net/?retryWrites=true&w=majority&appName=DB-1"
const DbName = "Ed_Fisic"
const CollectionName = "Ed_Fisic"

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
