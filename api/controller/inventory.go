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

func EditInventoryHandler(c *fiber.Ctx)error{
	inventoryID := c.Params("id")
	inventory,err := model.EditInventory(c, inventoryID)
	if err != nil{
		utilities.ShowError(c,"failed to edit inventory",0,map[string][]string{
			"error":{err.Error()},
		})
	}
	return utilities.ShowSuccess(c,"successfully updated inventory",1,inventory)
}

func DeleteInventoryHandler(c *fiber.Ctx)error{
	inventoryID := c.Params("id")
	err := model.DeleteInventory(c,inventoryID)
	if err != nil{
		return utilities.ShowError(c,"failed to delete record",0, map[string][]string{
			"error":{err.Error()},
		})
	}
	return utilities.ShowMessage(c,"record deleted successfully",1)
}

func GetInventoryHandler(c *fiber.Ctx)error{
	inventory, err := model.GetInventory(c)
	if err != nil{
		return utilities.ShowError(c,"failed to get inventory record",0, map[string][]string{
			"errors":{err.Error()},
		})
	}
	return utilities.ShowSuccess(c,"successfully retrieved inventory records",1,inventory)
}