package models

import "time"

type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Permission int
}

// Session represents a user session
type Session struct {
	Username string
	Expiry   time.Time
}
