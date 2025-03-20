package reception

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/dancankarani/medicare/api/controller/user"
	"github.com/dancankarani/medicare/api/model"
	"github.com/gofiber/fiber/v2"
)
func SetReceptionRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/reception")
	auth.Post("/login",user.Login)

	receptionGroup := auth.Group("/",user.JWTMiddleware)
	//protected routes
	receptionGroup.Post("/",model.CreatePatient)
	receptionGroup.Get("/:id",controller.GetPatientBillsHanlder)
	receptionGroup.Patch("/",controller.UpdatePatientHandler)
}