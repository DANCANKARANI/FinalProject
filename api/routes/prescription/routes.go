package prescription

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/dancankarani/medicare/api/controller/user"
	"github.com/gofiber/fiber/v2"
)

func SetPrescriptionRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/prescription")
	auth.Post("/", controller.CreatePrescriptionHandler)
	presGroup := auth.Group("/",user.JWTMiddleware)
	presGroup.Get("/", controller.GetPrescriptionsHandler)
	presGroup.Get("/:id", controller.GetPrescriptionHandler)
	presGroup.Patch("/:id", controller.UpdatePrescriptionHandler)
	presGroup.Delete("/:id", controller.DeletePrescriptionHandler)
	//protected routes
}