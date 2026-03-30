package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GalleryImage represents an image saved to the user's gallery.
type GalleryImage struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Image     string             `json:"image" bson:"image"` // Base64 PNG
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
