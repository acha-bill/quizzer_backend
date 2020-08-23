package main

import (
	"github.com/acha-bill/quizzer_backend/plugins"
	"context"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	mongoURL := os.Getenv("MONGODB_URL")
	dbName := os.Getenv("DATABASE_NAME")

	clientOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database(dbName).Collection("testCollection")
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))


	// Routes
	for _, plugin := range plugins.Plugins {
		for _, handler := range plugin.Handlers() {
			path := plugin.Name() + "/" + handler.Path
			e.Add(handler.Method, path , handler.Handler)
		}
	}

	// Start server
	e.Logger.Fatal(e.Start(":8081"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
