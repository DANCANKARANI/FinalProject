package controller

import (
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
)

func CreatePrescriptionHandler(c *fiber.Ctx) error {
	prescription, err := model.CreatePrescription(c)
	if err != nil{
		return utilities.ShowError(c,"failed to create upload prescription",0,map[string][]string{
			"error":{err.Error()},
		})
	}
	return utilities.ShowSuccess(c,"successfully uploaded the prescription",1,prescription)
}

//get the prescriptions handler
func GetPrescriptionsHandler(c *fiber.Ctx)error{
	prescriptions,err := model.GetPrescriptions(c)
	if err != nil{
		return utilities.ShowError(c,"failed to get prescriptions",0,map[string][]string{
			"error": {err.Error()},
		})
	}
	return utilities.ShowSuccess(c,"successfully retrieved prescriptions",1,prescriptions)
}

//get prescription handler
func GetPrescriptionHandler(c *fiber.Ctx)error{
	id := c.Params("id")
	prescription,err := model.GetPrescription(c,id)
	if err != nil{
		return utilities.ShowError(c,"failed to get prescription",0,map[string][]string{
			"error":{err.Error()},
		})
	}
	return utilities.ShowSuccess(c,"successfully retrieved prescription",1,prescription)
}

//updated prescription handler
func UpdatePrescriptionHandler(c *fiber.Ctx)error{
	id := c.Params("id")
	prescription,err := model.UpdatePrescription(c,id)
	if err != nil{
		return utilities.ShowError(c,"failed to update prescription",0,map[string][]string{
			"error":{err.Error()},
		})
	}
	return utilities.ShowSuccess(c,"prescriptions updated successfully",1,prescription)
}

//delete prescription handler
func DeletePrescriptionHandler(c *fiber.Ctx)error{
	id := c.Params("id")
	err := model.DeletePrescription(c,id)
	if err != nil{
		return utilities.ShowError(c,"failed to delete prescription",0,map[string][]string{
			"error":{err.Error()},
		})
	}
	return utilities.ShowMessage(c,"prescription deleted successfully",1)
}