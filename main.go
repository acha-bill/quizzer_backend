package main

import (
	"github.com/acha-bill/quizzer_backend/packages/mongodb"
	"github.com/acha-bill/quizzer_backend/packages/server"
	"github.com/acha-bill/quizzer_backend/packages/socketserver"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

// @title Quizzer API
// @version 1.0
// @description API for quizzer
// @termsOfService http://swagger.io/terms/

// @contact.name Acha Bill
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	_, err = mongodb.Connect()
	if err != nil {
		log.Fatal(err)
	}

	e := server.Instance()
	e.Static("/", "./public")
	e.GET("/ws", socketserver.Listen)
	e.Logger.Fatal(e.Start(":8081"))
}
