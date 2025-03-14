package patient

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/dancankarani/medicare/api/model"
	"github.com/gofiber/fiber/v2"
)

func SetPatientRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/patient")
	auth.Post("/",controller.CreatePatientHandler)
	auth.Put("/:id",controller.UpdatePatientHandler)
	auth.Delete("/:id",model.DeletePatient)
	auth.Get("/",model.GetPatients)

	//refer patient
	auth.Post("/:id",controller.ReferPatientHandler)
}