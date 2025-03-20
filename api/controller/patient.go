package controller

import (
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreatePatientHandler(c *fiber.Ctx)error{
	err := model.CreatePatient(c)
	if err != nil{
		utilities.ShowError(c,"failed to add patient",1, map[string][]string{
			"errors":{err.Error()},
		})
	}
	return utilities.ShowSuccess(c,"patient added successfully",0,nil)
}

//edit patient handler
func UpdatePatientHandler(c *fiber.Ctx)error{
	patient,err := model.UpdatePatient(c)
	if err != nil{
		return utilities.ShowError(c,"failed to update patient",1,map[string][]string{"error":{err.Error()}})
	}
	return utilities.ShowSuccess(c,"patient updated successfully",0,patient)
}

func GetPatientBillsHanlder(c *fiber.Ctx)error{
	patient_id,_ := uuid.Parse(c.Params("id"))
	bill, err := model.GetPatientBills(patient_id)
	if err != nil{
		return utilities.ShowError(c,"failed to get bills",1, map[string][]string{"errors":{err.Error()}})
	}
	return utilities.ShowSuccess(c,"bills retrieved successfully",0,bill)
}