package role

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/gofiber/fiber/v2"
)

func SetRoleRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/roles")
	//add middleware
	auth.Post("/", controller.CreateRoleHandler)
	auth.Patch("/:id",controller.UpdateRoleHandler)
	auth.Get("/",controller.GetAllRolesHandler)
	auth.Get("/:id",controller.GetRoleHandler)
	auth.Delete("/:id",controller.DeleteRoleHandler)
	//protected routes
}