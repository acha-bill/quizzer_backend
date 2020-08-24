package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Question represents a question
type Question struct {
	ID            primitive.ObjectID `bson:"id"`
	Question      string             `bson:"question"`
	Answers       []string           `bson:"answers"`
	CorrectAnswer string             `bson:"correctAnswer"`
	CreatedAt     time.Time          `bson:"created_At"`
	UpdatedAt     time.Time          `bson:"updated_At"`
}
