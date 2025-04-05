package routes

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/")
	UserRoutes(api)
	TaskRoutes(api)
	AuthRoutes(api)
}
