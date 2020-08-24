package question

import (
	"context"
	"errors"

	"github.com/acha-bill/quizzer_backend/models"
	"github.com/acha-bill/quizzer_backend/packages/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName = "questions"
)

var (
	ctx                  = context.TODO()
	ErrNoQuestionDeleted = errors.New("no questions were deleted")
)

func collection() *mongo.Collection {
	db, _ := mongodb.Database()
	return db.Collection(collectionName)
}

func FindAll() (questions []*models.Question, err error) {
	// passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}
	questions, err = filterQuestions(filter)
	return
}

func Find(filter interface{}) (questions []*models.Question, err error) {
	questions, err = filterQuestions(filter)
	return
}

func Create(question models.Question) (created *models.Question, err error) {
	res, err := collection().InsertOne(ctx, question)
	if err != nil {
		return nil, err
	}
	question.ID = res.InsertedID.(primitive.ObjectID)
	created = &question
	return
}

func UpdateById(id string, question models.Question) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	b, _ := bson.Marshal(&question)
	update := bson.D{primitive.E{Key: "$set", Value: b}}
	updated := &models.Question{}
	_ = collection().FindOneAndUpdate(ctx, filter, update).Decode(updated)
}

func DeleteById(id string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	res, err := collection().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrNoQuestionDeleted
	}

	return nil
}

func filterQuestions(filter interface{}) ([]*models.Question, error) {
	var questions []*models.Question

	cur, err := collection().Find(ctx, filter)
	if err != nil {
		return questions, err
	}

	for cur.Next(ctx) {
		var u models.Question
		err := cur.Decode(&u)
		if err != nil {
			return questions, err
		}

		questions = append(questions, &u)
	}

	if err := cur.Err(); err != nil {
		return questions, err
	}

	// once exhausted, close the cursor
	_ = cur.Close(ctx)

	if len(questions) == 0 {
		return questions, nil
	}

	return questions, nil
}
