package labtechnician

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/dancankarani/medicare/api/controller/user"
	"github.com/gofiber/fiber/v2"
)

func SetLabTechnicianRoutes(app *fiber.App){
	auth := app.Group("/api/v1/technician")
	auth.Post("/login",user.Login)

	technicianGroup := auth.Group("/",user.JWTMiddleware)
	technicianGroup.Post("/",controller.UploadLabTestHandler)
	technicianGroup.Get("/",controller.GetAllLabTestHandler)
	technicianGroup.Get("/:id",controller.GetLabTestByIdHandler)
}