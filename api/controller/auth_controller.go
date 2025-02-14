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
    pharmacist,err := model.CreateUser(c, user, true, false, "")
    if err != nil {
        return utilities.ShowError(c, "Failed to create pharmacist", fiber.StatusInternalServerError, map[string][]string{
            "error": {err.Error()},
        })
    }

    return utilities.ShowSuccess(c, "Pharmacist created successfully", 1, pharmacist)
}

/*
creates doctor
*/
func CreateDoctorHandler(c *fiber.Ctx) error {
    var user model.User
    // Parse request body
    if err := c.BodyParser(&user); err != nil {
        return utilities.ShowError(c, "Invalid request body", 0, map[string][]string{
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

    // Validate specialty
    var requestData struct {
        Specialty string `json:"specialty"`
    }

    if err := c.BodyParser(&requestData); err != nil {
        return utilities.ShowError(c, "Invalid specialty data", 0, map[string][]string{
            "body": {"Unable to parse specialty JSON"},
        })
    }

    if requestData.Specialty == "" {
        errors["specialty"] = append(errors["specialty"], "Doctor's specialty is required")
    }

    user.Password = utilities.PasswordGenerator()

    // Return validation errors if any
    if len(errors) > 0 {
        return utilities.ShowError(c, "Validation failed", 0, errors)
    }

    // Create user and assign as a doctor
    doctor,err := model.CreateUser(c, user, false, true, requestData.Specialty)
    if err != nil {
        return utilities.ShowError(c, "Failed to create doctor", 0, map[string][]string{
            "error": {err.Error()},
        })
    }

    return utilities.ShowSuccess(c, "Doctor created successfully", 1, doctor)
}
/*func DeleteUserHandler(c *fiber.Ctx) error {
    userID := c.Params("id")

    // Call DeleteUser function
    err := model.DeleteUser(c, userID)
    if err != nil {
        if strings.Contains(err.Error(), "User not found") {
            return utilities.ShowError(c, "User not found", 0, map[string][]string{
                "error": {err.Error()},
            })
        }
        return utilities.ShowError(c, "Failed to delete user", 0, map[string][]string{
            "error": {err.Error()},
        })
    }

    return utilities.ShowMessage(c, "User deleted successfully", 1)
}
*/

func GetUserHandler(c *fiber.Ctx)error{
    userID := c.Params("id")
    user, err := model.GetOneUser(c,userID)
    if err != nil{
        return utilities.ShowError(c,"failed to get user",0,map[string][]string{
            "errors":{err.Error()},
        })
    }
    return utilities.ShowSuccess(c,"successfully retrieved user",1,user)
}