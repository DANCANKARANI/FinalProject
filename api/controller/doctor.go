package controller

import (
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func EditDoctorHandler(c *fiber.Ctx) error {
    type Request struct {
        FullName  string `json:"full_name"`
        Email     string `json:"email"`
        Username  string `json:"username"`
        Password  string `json:"password,omitempty"`
        PhoneNumber string `json:"phone_number"`
        Specialty string `json:"specialty"`
    }

    userID, err := uuid.Parse(c.Params("id")) // Get user_id from URL params
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
    validPhone, err := utilities.ValidatePhoneNumber(req.PhoneNumber, "KE") // Use country code
    if err != nil {
        return utilities.ShowError(c, err.Error(), 0, nil)
    }

    updatedUser := model.User{
        FullName: req.FullName,
        Email:    *validEmail,
        Username: req.Username,
        PhoneNumber: validPhone,
        Password: req.Password,
    }

    updated_doctor,err := model.EditUser(c, userID, updatedUser, false, true, req.Specialty);
	if err != nil {
        return utilities.ShowError(c, err.Error(), 0, nil)
    }

    return utilities.ShowSuccess(c, "Doctor updated successfully", 1, updated_doctor)
}

/*Gets all the doctors*/
func GetDoctorsHandler(c *fiber.Ctx) error {
    users, err := model.GetDoctors(c)
    if err != nil {
        return utilities.ShowError(c, err.Error(), 0, nil)
    }
    return utilities.ShowSuccess(c, "Doctors fetched successfully", 1, users)
}

/*get all users*/
func GetAllUsersHandler(c *fiber.Ctx) error {
    users, err := model.GetAllUsers(c)
    if err != nil {
        return utilities.ShowError(c, err.Error(), fiber.StatusInternalServerError, nil)
    }
    return utilities.ShowSuccess(c, "Users fetched successfully", fiber.StatusOK, users)
}