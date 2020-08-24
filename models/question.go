package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Question represents a question
type Question struct {
	ID            primitive.ObjectID `bson:"_id"`
	Question      string             `bson:"question"`
	Answers       []string           `bson:"answers"`
	CorrectAnswer string             `bson:"correctAnswer"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}
