package controller

import (
	"log"

	"github.com/dancankarani/medicare/internal/app/model"
	"github.com/dancankarani/medicare/pkg/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateRoleHandler(c *fiber.Ctx) error {
	role := new(model.Role)

	if err := c.BodyParser(&role); err != nil{
		log.Println("error parsing request body:",err.Error())
		return utilities.ShowError(c,"invalid request",1,nil)
	}

	//call create error
	role,err := model.CreateRole(*role)
	if err != nil{
		return utilities.ShowError(c,err.Error(),1,nil)
	}

	return utilities.ShowSuccess(c,"role created successfully",0,role)
}

//update role handler
func UpdateRoleHandler(c *fiber.Ctx)error{
	role := new(model.Role)
	role_id,_ := uuid.Parse(c.Params("id"))
	if err := c.BodyParser(&role); err != nil{
		log.Println("error parsing request body:",err.Error())
		return utilities.ShowError(c,"invalid request",1,nil)
	}

	role, err := model.UpdateRole(role_id,*role)
	if err != nil{
		return utilities.ShowError(c,err.Error(),1,nil)
	}

	return utilities.ShowSuccess(c,"role updated successfully",0,role)
}

//get role by id
func GetRoleHandler(c *fiber.Ctx)error{
	role_id,_ := uuid.Parse(c.Params("id"))
	role,err := model.GetRole(role_id)
	if err != nil{
		return utilities.ShowError(c,err.Error(),1,nil)
	}
	
	return utilities.ShowSuccess(c,"Successfully retrieved role",0,role)
}

//get all roles handler
func GetAllRolesHandler(c *fiber.Ctx)error{
	roles, err := model.GetAllRoles()
	if err != nil{
		return utilities.ShowError(c,err.Error(),1,nil)
	}

	return utilities.ShowSuccess(c,"Successfully retrieved all roles",0,roles)
}

//delete role handler
func DeleteRoleHandler(c *fiber.Ctx)error{
	role_id:= uuid.MustParse(c.Params("id"))
	err := model.DeleteRole(role_id)
	if err != nil{
		return utilities.ShowError(c,err.Error(),1,nil)
	}

	return utilities.ShowMessage(c,"role deleted successfully",0)
}