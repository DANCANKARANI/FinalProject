package user

import (
	"github.com/dancankarani/medicare/internal/controller"
	"github.com/gofiber/fiber/v2"
)
func SetUserRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/users")
	auth.Post("/",controller.CreateUserHandler)
	//protected routes
}