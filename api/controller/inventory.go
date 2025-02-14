package controller

import (
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
)

func CreateInventoryHandler(c *fiber.Ctx)error{
	inventory,err := model.CreateInventory(c)
	if err != nil {
		return utilities.ShowError(c,"failed to add inventory",0,map[string][]string{
			"error":{err.Error()},
		})
	}
	return utilities.ShowSuccess(c,"inventory added successfully",1, inventory)
}
