package models

import "time"

// UserStatus defines permission levels for users
type UserStatus int
const (
    StatusPending          UserStatus = 0
    StatusReaderTimeBased  UserStatus = 1
    StatusReaderAmountBased UserStatus = 2
    StatusEditorTimeBased  UserStatus = 3
    StatusEditorAmountBased UserStatus = 4
    StatusEditorUnlimited  UserStatus = 5
    StatusAdmin            UserStatus = 6
)

// String returns the string representation of a UserStatus
func (s UserStatus) String() string {
    switch s {
    case StatusPending:
        return "Pending"
    case StatusReaderTimeBased:
        return "Reader (Time Based)"
    case StatusReaderAmountBased:
        return "Reader (Amount Based)"
    case StatusEditorTimeBased:
        return "Editor (Time Based)"
    case StatusEditorAmountBased:
        return "Editor (Amount Based)"
    case StatusEditorUnlimited:
        return "Editor (Unlimited)"
    case StatusAdmin:
        return "Administrator"
    default:
        return "Unknown"
    }
}

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
