package category

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
	collectionName = "categories"
)

var (
	ctx                    = context.TODO()
	ErrNoCategoriesDeleted = errors.New("no categories were deleted")
)

func collection() *mongo.Collection {
	db, _ := mongodb.Database()
	return db.Collection(collectionName)
}

func FindAll() (categories []*models.Category, err error) {
	// passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}
	categories, err = filterCategories(filter)
	return
}

func Find(filter interface{}) (categories []*models.Category, err error) {
	categories, err = filterCategories(filter)
	return
}

func Create(category models.Category) (created *models.Category, err error) {
	res, err := collection().InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}
	category.ID = res.InsertedID.(primitive.ObjectID)
	created = &category
	return
}

func FindById(id string) (category *models.Category, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	categories, err := filterCategories(filter)
	if err != nil {
		return
	}
	if len(categories) == 0 {
		category = nil
	} else {
		category = categories[0]
	}
	return
}

func UpdateById(id string, category models.Category) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	b, _ := bson.Marshal(&category)
	update := bson.D{primitive.E{Key: "$set", Value: b}}
	updated := &models.Category{}
	return collection().FindOneAndUpdate(ctx, filter, update).Decode(updated)
}

func DeleteById(id string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	res, err := collection().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrNoCategoriesDeleted
	}

	return nil
}

func filterCategories(filter interface{}) ([]*models.Category, error) {
	var categories []*models.Category

	cur, err := collection().Find(ctx, filter)
	if err != nil {
		return categories, err
	}

	for cur.Next(ctx) {
		var u models.Category
		err := cur.Decode(&u)
		if err != nil {
			return categories, err
		}

		categories = append(categories, &u)
	}

	if err := cur.Err(); err != nil {
		return categories, err
	}

	// once exhausted, close the cursor
	_ = cur.Close(ctx)

	if len(categories) == 0 {
		return categories, nil
	}

	return categories, nil
}
