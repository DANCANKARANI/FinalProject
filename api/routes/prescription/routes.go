package prescription

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/gofiber/fiber/v2"
)

func SetPrescriptionRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/prescription")
	auth.Post("/", controller.CreatePrescriptionHandler)
	auth.Get("/", controller.GetPrescriptionsHandler)
	auth.Get("/:id", controller.GetPrescriptionHandler)
	auth.Patch("/:id", controller.UpdatePrescriptionHandler)
	auth.Delete("/:id", controller.DeletePrescriptionHandler)
	//protected routes
}