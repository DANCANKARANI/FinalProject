package payments

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/gofiber/fiber/v2"
)

func SetPaymentsRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/payments")
	auth.Post("/", controller.MakePayments)
	auth.Post("/callback",controller.HandleCallback)
	//protected routes
}