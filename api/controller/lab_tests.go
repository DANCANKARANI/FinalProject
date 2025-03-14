package controller

import (
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadLabTestHandler(c *fiber.Ctx)error{
	LabTest, err := model.UploadLabTest(c)
	if err != nil{
		return utilities.ShowError(c,"failed to upload lab test",1,map[string][]string{"errors":{err.Error()}})
	}
	return utilities.ShowSuccess(c,"lab test uploaded",0,LabTest)
}

func GetAllLabTestHandler(c *fiber.Ctx)error{
	labTests, err := model.GetAllLabTest(c)
	if err != nil{
		return utilities.ShowError(c,"failed to get lab tests",1,map[string][]string{"errors":{err.Error()}})
	}
	return utilities.ShowSuccess(c,"successfully retrieved all lab tests",0,labTests)
}

func GetLabTestByIdHandler(c *fiber.Ctx)error{
	lab_test_id,_:= uuid.Parse(c.Params("id"))
	labTest, err := model.GetLabTestByID(c,lab_test_id)
	if err != nil{
		return utilities.ShowError(c,"failed to get lab test",1,map[string][]string{"errors":{err.Error()}})
	}
	return utilities.ShowSuccess(c,"successfully retrieved tlab test",0,labTest)
}


