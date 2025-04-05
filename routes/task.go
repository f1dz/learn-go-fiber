package routes

import (
	"fiber-api/controllers"
	"fiber-api/middlewares"

	"github.com/gofiber/fiber/v3"
)

func TaskRoutes(router fiber.Router) {
	api := router.Group("/task")
	api.Use(middlewares.JWTMiddleware())
	api.Post("/", controllers.CreateTask)
	api.Get("/", controllers.GetTasks)
	api.Get("/user", controllers.GetUserTasks)
	api.Get("/:id", controllers.GetTask)
	api.Put("/:id", controllers.UpdateTask)
	api.Delete("/:id", controllers.DeleteTask)
	api.Get("/user/:id", controllers.GetUserTasks)
	api.Get("/status/:status", controllers.GetTaskByStatus)
}
