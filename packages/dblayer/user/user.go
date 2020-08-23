package user

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
	collectionName = "users"
)
var (
	ctx = context.TODO()
	ErrNoUserDeleted = errors.New("no users were deleted")
)

func collection () *mongo.Collection{
	db, _ := mongodb.Database()
	return db.Collection(collectionName)
}

func FindAll() (users []*models.User, err error) {
		// passing bson.D{{}} matches all documents in the collection
		filter := bson.D{{}}
		users, err = filterUsers(filter)
		return
}

func Find(filter interface{}) (users []*models.User,  err error) {
	users, err = filterUsers(filter)
	return
}

func Create(user models.User) (created *models.User, err error) {
	res, err := collection().InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	created = &user
	return
}

func UpdateById(id string, user models.User) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	b,_ := bson.Marshal(&user)
	update := bson.D{primitive.E{Key: "$set", Value: b}}
	updated := &models.User{}
	_ = collection().FindOneAndUpdate(ctx, filter, update).Decode(updated)
}

func DeleteById(id string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	res, err := collection().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrNoUserDeleted
	}

	return nil
}

func filterUsers(filter interface{}) ([]*models.User, error) {
	var users []*models.User

	cur, err := collection().Find(ctx, filter)
	if err != nil {
		return users, err
	}

	for cur.Next(ctx) {
		var u models.User
		err := cur.Decode(&u)
		if err != nil {
			return users, err
		}

		users = append(users, &u)
	}

	if err := cur.Err(); err != nil {
		return users, err
	}

	// once exhausted, close the cursor
	_ = cur.Close(ctx)

	if len(users) == 0 {
		return users, nil
	}

	return users, nil
}