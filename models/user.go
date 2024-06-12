package models

import "time"

type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Permission int  `json:"permission" bson:"permission"`
}

// Session represents a user session
type Session struct {
	Username string
	Expiry   time.Time
}
