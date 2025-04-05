package main

import (
	"fiber-api/config"
	"fiber-api/routes"
	"fiber-api/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	utils.InitValidator()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	config.ConnectDB()
	app := fiber.New()
	routes.RegisterRoutes(app)

	port := os.Getenv("APP_PORT")
	app.Listen(":" + port)
}
