package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatMessage represents a single chat message stored in MongoDB.
type ChatMessage struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChatID    string             `json:"chatId" bson:"chatId"`
	Username  string             `json:"username" bson:"username"`
	Type      string             `json:"type" bson:"type"`           // "user" or "assistant"
	Content   string             `json:"content" bson:"content"`
	Image     string             `json:"image,omitempty" bson:"image,omitempty"` // Base64 PNG, nullable
	HasImage  bool               `json:"hasImage" bson:"hasImage"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
}
