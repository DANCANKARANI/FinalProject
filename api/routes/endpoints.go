package routes

import (
	"github.com/dancankarani/medicare/api/routes/admin/doctor"
	"github.com/dancankarani/medicare/api/routes/admin/pharmacist"
	"github.com/dancankarani/medicare/api/routes/inventory"
	labtechnician "github.com/dancankarani/medicare/api/routes/lab_technician"
	"github.com/dancankarani/medicare/api/routes/medicine"
	"github.com/dancankarani/medicare/api/routes/patient"
	"github.com/dancankarani/medicare/api/routes/prescription"
	"github.com/dancankarani/medicare/api/routes/reception"
	"github.com/dancankarani/medicare/api/routes/role"
	"github.com/dancankarani/medicare/api/routes/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func RegisterEndpoints() {
	app := fiber.New()
    app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins, change this to specific origins in production
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization", 
	}))
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
    //
    app.Listen(":8000")
}