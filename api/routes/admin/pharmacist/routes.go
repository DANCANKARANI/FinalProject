package pharmacist

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/dancankarani/medicare/api/model"
	"github.com/gofiber/fiber/v2"
)

func SetPharmacistRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/admin/pharmacist")
	auth.Post("/",controller.CreatePharmacistHandler)
	auth.Put("/:id",controller.EditPharmacistHandler)
	auth.Get("/",controller.GetPharmacistsHandler)
	auth.Delete("/:id",model.DeleteUser)
}