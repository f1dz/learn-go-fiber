package routes

import (
	"fiber-api/controllers"
	"fiber-api/middlewares"

	"github.com/gofiber/fiber/v3"
)

func UserRoutes(router fiber.Router) {
	api := router.Group("/user")
	api.Post("/register", controllers.Register)
	api.Use(middlewares.JWTMiddleware())
	api.Get("/", controllers.GetUsers)
	api.Get("/me", controllers.GetProfile)
	api.Get("/:id", controllers.GetUser)
}
