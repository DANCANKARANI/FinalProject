package inventory

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/gofiber/fiber/v2"
)

func SetInventoryRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/inventory")
	//inventory routes
	auth.Post("/",controller.CreateInventoryHandler)
	auth.Put("/:id",controller.EditInventoryHandler)
	auth.Delete("/:id",controller.DeleteInventoryHandler)
	auth.Get("/",controller.GetInventoryHandler)
}