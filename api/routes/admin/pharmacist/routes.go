package pharmacist

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/gofiber/fiber/v2"
)

func SetPahrmacistRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/admin/pharmacist")
	auth.Post("/",controller.CreatePharmacistHandler)

}