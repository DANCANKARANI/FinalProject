package user

import (
	"log"
	"time"

	"github.com/dancankarani/medicare/api/middleware"
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
)
type ResponseUser struct{
	FullName string 	`json:"full_name"`
	PhoneNumber string 	`json:"phone_number"`
	Email string 		`json:"email"`
}
type loginResponse struct {
	Token string `json:"token"`
}

func Login(c *fiber.Ctx) error {
    user := &model.User{}
    if err := c.BodyParser(&user); err != nil {
        return utilities.ShowError(c, "failed to login", 1, map[string][]string{
            "errors": {err.Error()},
        })
    }

    // Check if user exists
    _,existingUser,err := model.UserExist(c,user.Username)
    if err != nil{
        log.Println("err:"+err.Error())
        return utilities.ShowError(c,"login failed",1, map[string][]string{"errors":{err.Error()}})
    }

    // Compare password
    pass := existingUser.Password
    err = utilities.CompareHashAndPassowrd(pass, user.Password)
    if err != nil {
        return utilities.ShowError(c, err.Error(), 1, map[string][]string{
            "errors": {err.Error()},
        })
    }

    // Generate token
    exp := time.Hour * 24
    tokenString, err := middleware.GenerateToken(middleware.Claims{
        UserID:   &existingUser.ID,
        Role:     existingUser.Role,
        FullName: existingUser.FullName,
    }, exp)
    if err != nil {
        return utilities.ShowError(c, err.Error(), 1, map[string][]string{"errors": {err.Error()}})
    }

    domain := c.Hostname()
    if domain == "localhost" {
        domain = "localhost:3000"
    } else {
        domain = "ehospital-ashy.vercel.app"
    }
    
    c.Cookie(&fiber.Cookie{
        Name:     "Authorization",
        Value:    tokenString,
        Expires:  time.Now().Add(24 * time.Hour),
        HTTPOnly: true, // Prevent client-side JavaScript access
        Secure:   true, // Enable in production
        SameSite: "Lax",
        Domain:   domain, // Dynamically set based on request
        Path:     "/",
    })

    // Return success response
    responseUser := loginResponse{
        Token: tokenString,
    }
    return utilities.ShowSuccess(c, "successfully logged in", fiber.StatusOK, responseUser)
}

//logut user
func Logout(c *fiber.Ctx) error {
	err :=LogoutService(c,"normal")
	if err != nil{
		return utilities.ShowError(c,err.Error(),1,map[string][]string{"errors":{err.Error()}})
	}
	return  utilities.ShowMessage(c,"user logged out successfully",fiber.StatusOK)
}