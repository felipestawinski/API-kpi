package models

import "time"

type User struct {
	Email       string `json:"email" bson:"email"`
	Password    string `json:"password" bson:"password"`
	Username    string `json:"username" bson:"username"`
	Institution string `json:"institution" bson:"institution"`
	Role        string `json:"role" bson:"role"`
	Permission  int    `json:"permission" bson:"permission"`
	ID          string `json:"id,omitempty" bson:"_id,omitempty"`
	Files []    string `json:"files,omitempty" bson:"files,omitempty"`
}
// Session represents a user session
type Session struct {
	Username string
	Expiry   time.Time
}
