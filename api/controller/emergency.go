package controller

import (
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
)

func GetEmergemcyCasesHandler(c *fiber.Ctx) error {
	patient, err := model.GetEmergencyCases(c)
	if err != nil{
		return utilities.ShowError(c,"failed to get emergency cases",1,map[string][]string{"error":{err.Error()}})
	}
	return utilities.ShowSuccess(c,"emergency cases retrieved successfully",0,patient)
}