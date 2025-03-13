package user

import (
	"time"

	"github.com/dancankarani/medicare/api/middleware"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
)

func LogoutService(c *fiber.Ctx, user_type string) error {

	//get token string
	tokenString, err := utilities.GetJWTToken(c)
	if err != nil {
		return utilities.ShowError(c, err.Error(), 1,map[string][]string{"errors":{err.Error()}})
	}

	//invalidate token
	err = middleware.InvalidateToken(tokenString)
	if err != nil {
		return utilities.ShowError(c, "failed to invalidate the token", 1,map[string][]string{"errors":{err.Error()}})
	}


	
	//set token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    "",
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
	})

	//response
	return nil
}