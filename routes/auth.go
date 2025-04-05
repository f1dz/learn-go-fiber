package routes

import (
	"fiber-api/controllers"

	"github.com/gofiber/fiber/v3"
)

func AuthRoutes(router fiber.Router) {
	api := router.Group("/auth")
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
}
