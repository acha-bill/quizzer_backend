package mongodb

import (
	"context"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	once   sync.Once
	client *mongo.Client
	dbName string
)

// Database connects to a mongodb database and returns the DB instance
func Database() (db *mongo.Database, err error) {
	client, err := Connect()
	if err != nil {
		return nil, err
	}
	db = client.Database(dbName)
	return
}

// Connect connects to mongodb and returns the client instance
func Connect() (c *mongo.Client, err error) {
	once.Do(func() {
		client, err = connect()
	})
	return client, err
}

func connect() (c *mongo.Client, err error) {
	var ctx = context.TODO()
	mongoURL := os.Getenv("MONGODB_URL")
	dbName = os.Getenv("DATABASE_NAME")

	clientOptions := options.Client().ApplyURI(mongoURL)
	c, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = c.Ping(ctx, nil)
	if err != nil {
		return c, err
	}
	return
}
