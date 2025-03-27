package controller

import (
	"fmt"
	"log"

	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/api/services"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateDoctorHandler creates a new doctor
func CreateUserHandler(c *fiber.Ctx) error {
	var user model.User

	// Parse request body
	if err := c.BodyParser(&user); err != nil {
		return utilities.ShowError(c, "Invalid request body", fiber.StatusBadRequest, map[string][]string{
			"body": {"Unable to parse JSON"},
		})
	}
    //generate username
    Username,_ := services.GenerateRoleBasedUsername(user.Role)

    //assign username
    user.Username= Username


	// Validate required fields
	errors := make(map[string][]string)

	if user.FullName == "" {
		errors["full_name"] = append(errors["full_name"], "Full name is required")
	}

	if user.Email == "" {
		errors["email"] = append(errors["email"], "Email is required")
	} else {
		validatedEmail, err := utilities.ValidateEmail(user.Email)
		if err != nil {
			errors["email"] = append(errors["email"], err.Error())
		} else {
			user.Email = *validatedEmail // Store validated email
		}
	}

	if user.Username == "" {
		errors["username"] = append(errors["username"], "Username is required")
	}

	if user.PhoneNumber == "" {
		errors["phone_number"] = append(errors["phone_number"], "Phone number is required")
	} else {
		// Validate and format phone number
		formattedPhone, err := utilities.ValidatePhoneNumber(user.PhoneNumber, "KE")
		if err != nil {
			errors["phone_number"] = append(errors["phone_number"], err.Error())
		} else {
			user.PhoneNumber = formattedPhone // Store formatted number
		}
	}
    pass := utilities.PasswordGenerator()
    hashed_passsword, _ := utilities.HashPassword(pass)

	log.Println("username:" + user.Username + " password:" + pass)

	// Format the login details
	logins := fmt.Sprintf("Username: %s\nPassword: %s", user.Username, pass)
	
	// Send email
	utilities.SendEmail(user.Email, "Iruma Dispensary Login Details", logins)
	
	// Assign the hashed password to the user
	user.Password = hashed_passsword

	// Return validation errors if any
	if len(errors) > 0 {
		return utilities.ShowError(c, "Validation failed", fiber.StatusBadRequest, errors)
	}

	// Create user and assign as a doctor
	doctor, err := model.CreateUser(c, user)
	if err != nil {
		return utilities.ShowError(c, "Failed to create user", fiber.StatusInternalServerError, map[string][]string{
			"error": {err.Error()},
		})
	}

	return utilities.ShowSuccess(c, "Doctor created successfully", fiber.StatusOK, doctor)
}
//get user by role

// DeleteUserHandler deletes a user by ID
func DeleteUserHandler(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Call DeleteUser function
	err := model.DeleteUser(c, userID)
	if err != nil {
		return utilities.ShowError(c, "Failed to delete user", fiber.StatusInternalServerError, map[string][]string{
			"error": {err.Error()},
		})
	}

	return utilities.ShowMessage(c, "User deleted successfully", fiber.StatusOK)
}

// GetUserHandler retrieves a user by ID
func GetUserHandler(c *fiber.Ctx) error {
	userID,_ := uuid.Parse(c.Params("id"))

	// Fetch user by ID
	user, err := model.GetOneUser(c, userID)
	if err != nil {
		return utilities.ShowError(c, "Failed to get user", fiber.StatusNotFound, map[string][]string{
			"errors": {err.Error()},
		})
	}

	return utilities.ShowSuccess(c, "Successfully retrieved user", fiber.StatusOK, user)
}

//get user by role handler
func GetUsersByRoleHandler(c *fiber.Ctx)error{
    role := c.Params("role")
    user, err := model.GetUsersByRole(c,role)
    if err != nil{
        return utilities.ShowError(c,"failed to get user",1,map[string][]string{"errors":{err.Error()}})
    }
    return utilities.ShowSuccess(c,"successfully retrieved user by role",0,user)
}