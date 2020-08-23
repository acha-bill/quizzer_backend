package models

import (
	"time"
)

// User represents a user
type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"update_at"`
	ProfileURL string `bson:"profileUrl"`
}