package controller

import (
	"github.com/dancankarani/medicare/database"
	
	"github.com/gofiber/fiber/v2"
)
var db = database.ConnectDB()

func CreatePharmacistHandler(c *fiber.Ctx)error{

}