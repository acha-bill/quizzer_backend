package main

import (
	"fmt"
	"github.com/acha-bill/quizzer_backend/packages/mongodb"
	"github.com/acha-bill/quizzer_backend/packages/server"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"os"
)


func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	fmt.Println("0")
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
fmt.Println("1")
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	log.Info(os.Getenv("MONGODB_URL"))

	_, err = mongodb.Connect()
	if err != nil {
		log.Fatal(err)
	}
	e := server.Instance()
	e.Logger.Fatal(e.Start(":8081"))
}
