package patient

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/dancankarani/medicare/api/controller/user"
	"github.com/dancankarani/medicare/api/model"
	"github.com/gofiber/fiber/v2"
)

func SetPatientRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/patient")
	
	auth.Put("/:id",controller.UpdatePatientHandler)
	auth.Delete("/:id",model.DeletePatient)
	auth.Get("/",model.GetPatients)

	//refer patient
	patGroup := auth.Group("/",user.JWTMiddleware)
	patGroup.Post("/",controller.CreatePatientHandler)
	patGroup.Post("/:id",controller.ReferPatientHandler)
}