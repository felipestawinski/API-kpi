package models

import "time"

// UserStatus defines permission levels for users
type UserStatus int
const (
    StatusPending          UserStatus = 0
    StatusReaderTimeBased  UserStatus = 1
    StatusReaderAmountBased UserStatus = 2
    StatusReaderUnlimited  UserStatus = 3
    StatusEditorTimeBased  UserStatus = 4
    StatusEditorAmountBased UserStatus = 5
    StatusEditorUnlimited  UserStatus = 6
    StatusAdmin            UserStatus = 7
)

// String returns the string representation of a UserStatus
func (s UserStatus) String() string {
    switch s {
    case StatusPending:
        return "Pendente"
    case StatusReaderTimeBased:
        return "Leitor (Por Tempo)"
    case StatusReaderAmountBased:
        return "Leitor (Por Requisição)"
    case StatusReaderUnlimited:
        return "Leitor (Permanente)"
    case StatusEditorTimeBased:
        return "Editor (Por Tempo)"
    case StatusEditorAmountBased:
        return "Editor (Por Requisição)"
    case StatusEditorUnlimited:
        return "Editor (Permanente)"
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
    AccessTime  string `json:"accesstime" bson:"accesstime"`
    ReqAmount   int    `json:"reqamount" bson:"reqamount"`
	ID          string `json:"id,omitempty" bson:"_id,omitempty"`
	Files []    string `json:"files,omitempty" bson:"files,omitempty"`
	ProfilePicture string `json:"profilePicture,omitempty" bson:"profilePicture,omitempty"`
}
// Session represents a user session
type Session struct {
	Username string
	Expiry   time.Time
}
