package doctor

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/dancankarani/medicare/api/controller/user"
	"github.com/gofiber/fiber/v2"
)

func SetDoctorsRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/admin/doctor")
	auth.Post("/login",user.Login)
	auth.Post("/",controller.CreateUserHandler)
	doctorGroup := auth.Group("/",user.JWTMiddleware)
	doctorGroup.Put("/:id",controller.EditDoctorHandler)
	doctorGroup.Get("/users",controller.GetAllUsersHandler)
	doctorGroup.Delete("/:id",controller.DeleteUserHandler)
}