package database

import (
"context"
"time"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
)

const MongoURI = "mongodb://localhost:27017"
const DbName = "userdb"
const CollectionName = "users"
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