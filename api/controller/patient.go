package controller

import (
	"errors"
	"github.com/dancankarani/medicare/api/model"
	"github.com/gofiber/fiber/v2"
)

func CreatePatientHandler(c *fiber.Ctx)error{
	err := model.CreatePatient(c)
	if err != nil{
		return errors.New(err.Error())
	}
	return nil
}