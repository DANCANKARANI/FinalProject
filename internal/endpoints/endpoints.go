package endpoints

import (
	"github.com/dancankarani/medicare/internal/routes/role"
	"github.com/dancankarani/medicare/internal/routes/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func RegisterEndPoints() {
	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins, change this to specific origins in production
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization", 
	}))

	user.SetUserRoutes(app)
	role.SetRoleRoutes(app)
	//listening port
	app.Listen(":8000")
}