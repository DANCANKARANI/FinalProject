package controller

import (
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ReferPatientHandler(c *fiber.Ctx) error {
	patient_id,_:= uuid.Parse(c.Params("id"))
	referral, err := model.ReferPatient(c,patient_id)
	if err != nil{
		return utilities.ShowError(c,"failed to generate referral", 1,map[string][]string{"error":{err.Error()}})
	}
	return utilities.ShowSuccess(c,"patient referral generated successfully",0,referral)
}