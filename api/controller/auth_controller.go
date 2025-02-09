package controller

import (
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/utilities"

	"github.com/gofiber/fiber/v2"
)

func CreatePharmacistHandler(c *fiber.Ctx) error {
    var user model.User

    // Parse request body
    if err := c.BodyParser(&user); err != nil {
        return utilities.ShowError(c, "Invalid request body", fiber.StatusBadRequest, map[string][]string{
            "body": {"Unable to parse JSON"},
        })
    }

    // Validate required fields
    errors := make(map[string][]string)

    if user.FullName == "" {
        errors["full_name"] = append(errors["full_name"], "Full name is required")
    }

    if user.Email == "" {
        errors["email"] = append(errors["email"], "Email is required")
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

    user.Password = utilities.PasswordGenerator()

    // Return validation errors if any
    if len(errors) > 0 {
        return utilities.ShowError(c, "Validation failed", fiber.StatusBadRequest, errors)
    }

    // Create user and assign as a pharmacist
    err := model.CreateUser(c, user, true, false, "")
    if err != nil {
        return utilities.ShowError(c, "Failed to create pharmacist", fiber.StatusInternalServerError, map[string][]string{
            "error": {err.Error()},
        })
    }

    return utilities.ShowSuccess(c, "Pharmacist created successfully", 1, user)
}
