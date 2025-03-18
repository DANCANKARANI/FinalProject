package controller

import (
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func EditPharmacistHandler(c *fiber.Ctx) error {
    type Request struct {
        FullName    string `json:"full_name"`
        Email       string `json:"email"`
        Username    string `json:"username"`
        PhoneNumber string `json:"phone_number"`
        Password    string `json:"password,omitempty"`
    }

    // Parse user_id from URL params
    userID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return utilities.ShowError(c, "Invalid user ID", 0, nil)
    }

    var req Request
    if err := c.BodyParser(&req); err != nil {
        return utilities.ShowError(c, "Invalid request", 0, nil)
    }

    // Validate email
    validEmail, err := utilities.ValidateEmail(req.Email)
    if err != nil {
        return utilities.ShowError(c, err.Error(), 0, nil)
    }

    // Validate phone number
    validPhone, err := utilities.ValidatePhoneNumber(req.PhoneNumber, "KE") // Assuming Kenya (+254)
    if err != nil {
        return utilities.ShowError(c, err.Error(), 0, nil)
    }

    updatedUser := model.User{
        FullName:    req.FullName,
        Email:       *validEmail,
        Username:    req.Username,
        PhoneNumber: validPhone,
        Password:    req.Password,
    }

    // Call EditUser to update the pharmacist's details
    updated_pharmacist,err := model.EditUser(c, userID, updatedUser );
	if err != nil {
        return utilities.ShowError(c, err.Error(), 0, nil)
    }

    return utilities.ShowSuccess(c, "Pharmacist updated successfully", 1, updated_pharmacist)
}

func GetPharmacistsHandler(c *fiber.Ctx)error{
    role := "pharmacist"
    pharmacists,err := model.GetUsersByRole(c,role)
    if err != nil{
        return utilities.ShowError(c,"fai;ed to get pharmacists",1,map[string][]string{"errors":{err.Error()}})
    } 
    return utilities.ShowSuccess(c,"pharmacists retrueved successfully",0,pharmacists)
}