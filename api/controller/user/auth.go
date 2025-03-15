package user

import (
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

func Login(c *fiber.Ctx)error{
	user := &model.User{}
	if err := c.BodyParser(&user); err !=nil {
		return utilities.ShowError(c,"failed to login",1,map[string][]string{
			"errors":{err.Error()},
		})
	}

	//check of user exist
	userExist,existingUser,err:= model.UserExist(c,user.Username)
	if ! userExist && err != nil {
		return utilities.ShowError(c,"user does not exist",1,map[string][]string{
			"errors":{err.Error()},
		})
	}
	pass := existingUser.Password
	//compare password
	err = utilities.CompareHashAndPassowrd(pass,user.Password)
	if err !=nil{
		return utilities.ShowError(c,err.Error(),1,map[string][]string{
			"errors":{err.Error()},
		})		 
	}
	exp :=time.Hour*24
	//generating token
	tokenString,err := middleware.GenerateToken(middleware.Claims{UserID: &existingUser.ID,Role:existingUser.Role,FullName: existingUser.FullName},exp)
	if err != nil{
		return utilities.ShowError(c,err.Error(),1,map[string][]string{"errors":{err.Error()}})
	}
	//set token cookie 
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24), // Same duration as the token
		HTTPOnly: true, // Important for security, prevents JavaScript access
		Secure:   true, // Use secure cookies in production
		Path:     "/",  // Make the cookie available on all routes
	})
	response_user:=loginResponse{
		Token: tokenString,
	}

	return utilities.ShowSuccess(c,"successfully logged in",fiber.StatusOK,response_user)	
}

//logut user
func Logout(c *fiber.Ctx) error {
	err :=LogoutService(c,"normal")
	if err != nil{
		return utilities.ShowError(c,err.Error(),1,map[string][]string{"errors":{err.Error()}})
	}
	return  utilities.ShowMessage(c,"user logged out successfully",fiber.StatusOK)
}