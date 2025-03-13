package controller

import (
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
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