package medicine

import (
	"github.com/dancankarani/medicare/api/controller"
	"github.com/gofiber/fiber/v2"
)

func SetMedicineRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/medicine")
	auth.Post("/", controller.CreateMedicineHandler)
	auth.Get("/",controller.GetMedicinesHandler)
	auth.Get("/:id", controller.GetMedicineHandler)
	auth.Patch("/:id",controller.UpdateMedicineHandler)
	auth.Delete("/:id",controller.DeleteMedicineHandler)
	//protected routes
}