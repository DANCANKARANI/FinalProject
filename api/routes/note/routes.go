package note

import (
	"github.com/dancankarani/medicare/api/controller/user"
	"github.com/dancankarani/medicare/api/model"
	"github.com/gofiber/fiber/v2"
)

func SetNotesRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/notes")

	//refer patient
	noteGroup := auth.Group("/", user.JWTMiddleware)
	noteGroup.Post("/", model.CreateNote)
	noteGroup.Get("/",model.GetNote)
}