package controller

import (
	"time"

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
        DateOfBirth time.Time `json:"date_of_birth"`
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
   

    // Validate phone number
   
    updatedUser := model.User{
        FullName: req.FullName,
        Email:    req.Email,
        Username: req.Username,
        DateOfBirth: req.DateOfBirth,
        PhoneNumber: req.PhoneNumber,
        Password: req.Password,
    }

    updated_doctor,err := model.EditUser(c, userID, updatedUser);
	if err != nil {
        return utilities.ShowError(c, err.Error(), 0, nil)
    }

    return utilities.ShowSuccess(c, "Doctor updated successfully", 1, updated_doctor)
}

/*Gets all the doctors*/

/*get all users*/
func GetAllUsersHandler(c *fiber.Ctx) error {
    users, err := model.GetAllUsers(c)
    if err != nil {
        return utilities.ShowError(c, err.Error(), fiber.StatusInternalServerError, nil)
    }
    return utilities.ShowSuccess(c, "Users fetched successfully", fiber.StatusOK, users)
}

func GetUserById(c *fiber.Ctx)error{
    user_id,_ := model.GetAuthUserID(c)
    doctor, err := model.GetOneUser(c,user_id)
    if err != nil{
        return utilities.ShowError(c,"failed to get user",1,map[string][]string{"errors":{err.Error()}})
    }
    return utilities.ShowSuccess(c,"user retrieved successfully",0,doctor)
}

func GetDoctorsHandler(c *fiber.Ctx)error{
    role := "doctor"
    pharmacists,err := model.GetUsersByRole(c,role)
    if err != nil{
        return utilities.ShowError(c,"failed to get pharmacists",1,map[string][]string{"errors":{err.Error()}})
    } 
    return utilities.ShowSuccess(c,"pharmacists retrueved successfully",0,pharmacists)
}

func GetReceptionHandler(c *fiber.Ctx)error{
    role := "receptionist"
    pharmacists,err := model.GetUsersByRole(c,role)
    if err != nil{
        return utilities.ShowError(c,"failed to get receptionists",1,map[string][]string{"errors":{err.Error()}})
    } 
    return utilities.ShowSuccess(c,"receptionists retrieved successfully",0,pharmacists)
}

func GetTechnicianHandler(c *fiber.Ctx)error{
    role := "technician"
    pharmacists,err := model.GetUsersByRole(c,role)
    if err != nil{
        return utilities.ShowError(c,"failed to get lab technicians",1,map[string][]string{"errors":{err.Error()}})
    } 
    return utilities.ShowSuccess(c,"lab technicians retrieved successfully",0,pharmacists)
}