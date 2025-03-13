package controller

import (
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
)

//create medicine handler
func CreateMedicineHandler(c *fiber.Ctx) error {
medicine,err := model.CreateMedicine(c)
if err != nil{
	return utilities.ShowError(c,"failed to add medicine",0,map[string][]string{
		"error creating medicine":{err.Error()},
	})
}
return utilities.ShowSuccess(c,"successfully created a medicine",1,medicine)
}

//get all the medicines handler
func GetMedicinesHandler(c *fiber.Ctx)error{
	medicines, err := model.GetMedicines(c)
	if err != nil{
		return utilities.ShowError(c,"failed to get medicines",0,map[string][]string{
			"error":{err.Error()},
		})
	}
	return utilities.ShowSuccess(c,"successfully retrieved the medicines",1,medicines)
}

//get one medicine
func GetMedicineHandler(c *fiber.Ctx)error{
	id := c.Params("id")
	medicine, err := model.GetMedicine(c,id)
	if err != nil{
		return utilities.ShowError(c,"failed to get medicine",0,map[string][]string{
			"error":{err.Error()},
		})
	}
	return utilities.ShowSuccess(c,"medicine retrieved successfully",1,medicine)
}

//update medicine handler
func UpdateMedicineHandler(c *fiber.Ctx)error{
	id := c.Params("id")
	medicine, err := model.UpdateMedicine(c, id)
	if err != nil{
		return utilities.ShowError(c,"failed to update medicine",0,map[string][]string{
			"error":{err.Error()},
		})
	}
	return utilities.ShowSuccess(c,"medicine updated successfully", 1, medicine)
}

//delete medicine handler
func DeleteMedicineHandler(c *fiber.Ctx)error{
	id := c.Params("id")
	err := model.DeleteMedicine(c,id)
	if err != nil{
		return utilities.ShowError(c,"failed to delete medicine",0, map[string][]string{
			"error":{err.Error()},
		})
	}
	return utilities.ShowMessage(c,"successfully updated medicine",1)
}