package routes

import (
	"github.com/dancankarani/medicare/api/routes/admin"
	"github.com/dancankarani/medicare/api/routes/admin/doctor"
	"github.com/dancankarani/medicare/api/routes/admin/pharmacist"
	"github.com/dancankarani/medicare/api/routes/inventory"
	labtechnician "github.com/dancankarani/medicare/api/routes/lab_technician"
	"github.com/dancankarani/medicare/api/routes/medicine"
	"github.com/dancankarani/medicare/api/routes/note"
	"github.com/dancankarani/medicare/api/routes/patient"
	"github.com/dancankarani/medicare/api/routes/payments"
	"github.com/dancankarani/medicare/api/routes/prescription"
	"github.com/dancankarani/medicare/api/routes/reception"
	"github.com/dancankarani/medicare/api/routes/role"
	"github.com/dancankarani/medicare/api/routes/socket"
	"github.com/dancankarani/medicare/api/routes/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func RegisterEndpoints() {
	app := fiber.New()

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return origin == "http://localhost:3000" || origin == "https://ehospital-ashy.vercel.app"
		},
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Content-Type, Authorization, X-Requested-With, Accept, Origin, Access-Control-Request-Method, Access-Control-Request-Headers, X-Patient-ID",
		AllowCredentials: true,
	}))
	
	

	payments.SetPaymentsRoutes(app)
	note.SetNotesRoutes(app)
	// Register routes
	socket.SetWebSocketRoutes(app)
	doctor.SetDoctorsRoutes(app)
	user.SetUserRoutes(app)
	role.SetRoleRoutes(app)
	patient.SetPatientRoutes(app)
	pharmacist.SetPharmacistRoutes(app)
	inventory.SetInventoryRoutes(app)
	medicine.SetMedicineRoutes(app)
	prescription.SetPrescriptionRoutes(app)
	reception.SetReceptionRoutes(app)
	labtechnician.SetLabTechnicianRoutes(app)
	admin.SetAdminRoutes(app)

	// Start the server
	app.Listen(":8000")
}