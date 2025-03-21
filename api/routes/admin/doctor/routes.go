package doctor

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/dancankarani/medicare/api/controller/user"
	"github.com/gofiber/fiber/v2"
)

func SetDoctorsRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/admin/doctor")
	auth.Post("/login",user.Login)
	
	doctorGroup := auth.Group("/",user.JWTMiddleware)
	doctorGroup.Post("/",controller.CreateUserHandler)
	doctorGroup.Put("/:id",controller.EditDoctorHandler)
	doctorGroup.Get("/all",controller.GetDoctorsHandler)
	doctorGroup.Get("/users",controller.GetAllUsersHandler)
	doctorGroup.Get("/",controller.GetUserById)
	doctorGroup.Delete("/:id",controller.DeleteUserHandler)

	//getting the lab test
	
}