package model

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreatePrescription handles creating a new prescription// CreatePrescription handles creating a new prescription
func CreatePrescription(c *fiber.Ctx) (*Prescription, error) {
	// Define the request structure
	type PrescriptionRequest struct {
		PatientID             uuid.UUID `json:"patient_id"`
		MadicineName		string		`json:"medicine_name"`
		DoctorID              uuid.UUID `json:"doctor_id"`
		Diagnosis             string    `json:"diagnosis"`
		Dosage				string		`json:"dosage"`
		Instructions		string		`json:"instructions"`
		Status                string    `json:"status"`
	}

	// Parse the request body
	var req PrescriptionRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return nil, fmt.Errorf("invalid request data")
	}

	// Validate required fields
	if req.PatientID == uuid.Nil {
		return nil, fmt.Errorf("patient_id is required")
	}
	if req.DoctorID == uuid.Nil {
		return nil, fmt.Errorf("doctor_id is required")
	}
	if req.Diagnosis == "" {
		return nil, fmt.Errorf("diagnosis is required")
	}
	if req.MadicineName ==""{
		return nil, fmt.Errorf("medicine name is required")

	}

	// Create the prescription
	prescription := Prescription{
		ID:                 uuid.New(),
		PatientID:          req.PatientID,
		MedicineName: 		req.MadicineName,
		DoctorID:           req.DoctorID,
		Diagnosis:          req.Diagnosis,
		Dosage: req.Dosage,
		Instructions: req.Instructions,
		Status:             req.Status,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Save the prescription to the database
	if err := db.Create(&prescription).Error; err != nil {
		log.Printf("Error saving prescription: %v", err)
		return nil, fmt.Errorf("could not save prescription: %v", err)
	}

	// Preload the Doctor and Patient details
	if err := db.Preload("Doctor").Preload("Patient").First(&prescription, "id = ?", prescription.ID).Error; err != nil {
		log.Printf("Error preloading Doctor and Patient details: %v", err)
		return nil, fmt.Errorf("error preloading Doctor and Patient details: %v", err)
	}

	// Return the created prescription with preloaded Doctor and Patient details
	return &prescription, nil
}

// GetPrescriptions retrieves all prescriptions with optional filtering
func GetPrescriptions(c *fiber.Ctx) (*[]Prescription, error) {
	var prescriptions []Prescription

	// Initialize the query with preloading for PrescribedMedicines, Doctor, and Patient
	query := db.Preload("Doctor").Preload("Patient")

	// Apply filters based on query parameters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if patientID := c.Query("patient_id"); patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}
	if doctorID := c.Query("doctor_id"); doctorID != "" {
		query = query.Where("doctor_id = ?", doctorID)
	}

	// Execute the query
	if err := query.Find(&prescriptions).Error; err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve prescriptions"})
	}

	return &prescriptions, nil
}

// GetPrescription retrieves a single prescription by ID
func GetPrescription(c *fiber.Ctx,id string) (*Prescription,error) {
	var prescription Prescription

	if err := db.Preload("PrescribedMedicines").First(&prescription, "id = ?", id).Error; err != nil {
		return nil,c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Prescription not found"})
	}

	return &prescription, nil
}


// UpdatePrescription updates an existing prescription
func UpdatePrescription(c *fiber.Ctx, id string) (*Prescription, error) {
	var prescription Prescription

	// Fetch the prescription including prescribed medicines (if needed)
	if err := db.First(&prescription, "id = ?", id).Error; err != nil {
		return nil, c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Prescription not found"})
	}

	// Parse request body
	var updateData Prescription
	if err := c.BodyParser(&updateData); err != nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data"})
	}

	// Update fields only if they are provided in the request body

	if updateData.Diagnosis != "" {
		prescription.Diagnosis = updateData.Diagnosis
	}
	if updateData.Status != "" {
		prescription.Status = updateData.Status
	}
	if updateData.MedicineName != ""{
		prescription.MedicineName=updateData.MedicineName
	}

	

	// Update timestamp
	prescription.UpdatedAt = time.Now()

	// Save the updated prescription
	if err := db.Save(&prescription).Error; err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update prescription"})
	}

	// Return the updated prescription
	return &prescription, nil
}


// DeletePrescription deletes a prescription by ID
func DeletePrescription(c *fiber.Ctx,id string) error {
	var prescription Prescription

	if err := db.First(&prescription, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Prescription not found"})
	}

	if err := db.Delete(&prescription).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete prescription"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Prescription deleted successfully"})
}
