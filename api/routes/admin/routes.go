package admin

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/dancankarani/medicare/api/controller/user"
	"github.com/gofiber/fiber/v2"
)

func SetAdminRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/admin/")
	auth.Post("/login",user.Login)
	userGroup := auth.Group("/",user.JWTMiddleware)
	auth.Post("/",controller.CreateUserHandler)
	auth.Get("/technician",controller.GetTechnicianHandler)
	auth.Get("/reception",controller.GetReceptionHandler)
	userGroup.Get("/",controller.GetUsersByRoleHandler)
}