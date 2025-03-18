package pharmacist

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/dancankarani/medicare/api/controller/user"
	"github.com/gofiber/fiber/v2"
)

func SetPharmacistRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/admin/pharmacist")
	auth.Post("/login",user.Login)
	auth.Post("/",controller.CreateUserHandler)
	auth.Get("/",controller.GetPharmacistsHandler)
	auth.Put("/:id",controller.EditPharmacistHandler)
	auth.Delete("/:id",controller.DeleteUserHandler)
}