package doctor

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/dancankarani/medicare/api/model"
	"github.com/gofiber/fiber/v2"
)

func SetDoctorsRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/admin/doctor")
	auth.Post("/",controller.CreateDoctorHandler)
	auth.Put("/:id",controller.EditDoctorHandler)
	auth.Get("/",controller.GetDoctorsHandler)
	auth.Get("/users",controller.GetAllUsersHandler)
	auth.Delete("/:id",model.DeleteUser)
}