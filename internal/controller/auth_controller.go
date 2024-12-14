package controller

import (
    "log"
	"github.com/dancankarani/medicare/internal/app/database"
	"github.com/dancankarani/medicare/internal/app/model"
	"github.com/dancankarani/medicare/internal/app/services"
	"github.com/dancankarani/medicare/pkg/utilities"
	"github.com/gofiber/fiber/v2"
)
var db = database.ConnectDB()
func CreateUserHandler(c *fiber.Ctx) error {
    user := new(model.User)
    db.AutoMigrate(user)

    // Parse request body
    if err := c.BodyParser(user); err != nil {
        return utilities.ShowError(c, "Invalid request", 1, nil)
    }

    // Initialize an empty map for errors
    errors := make(map[string][]string)

    // Validate email
    if _, err := utilities.ValidateEmail(user.Email); err != nil {
        errors["email"] = append(errors["email"], "The email field must be a valid email address.")
    }

    // Validate phone number (Assuming 'KE' for Kenya)
    if _, err := utilities.ValidatePhoneNumber(user.PhoneNumber, "KE"); err != nil {
        errors["phone"] = append(errors["phone"], "The phone number must be a valid Kenyan number.")
    }

    // Check if there are any errors in the validation map
    if len(errors) > 0 {
        return utilities.ShowError(c, "Validation failed", 1, errors)
    }

	//check user existence
	mapStr,err := services.CheckUserRegistered(user.Email,user.PhoneNumber)
	if err != nil{
		return utilities.ShowError(c,err.Error(),1,mapStr)
	}

	
	username,_ := services.GenerateRoleBasedUsername("Doctor")
    user.Username=username
	user.Password=utilities.PasswordGenerator()
	log.Println(user.Username,user.Password)
    // Success message if no errors
    err = model.CreateUser(c, *user)
	if err != nil {
		return utilities.ShowError(c, err.Error(), 1, map[string][]string{})
	}

    return utilities.ShowMessage(c,"user created successfully",0)
}