package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User represents a user
type User struct {
	ID primitive.ObjectID `bson:"id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"update_at"`
	ProfileURL string `bson:"profileURL"`
	IsAdmin bool `bson:"isAdmin"`
}