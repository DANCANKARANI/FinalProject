package model

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreatePrescription handles creating a new prescription
func CreatePrescription(c *fiber.Ctx) (*Prescription, error) {
	// Define the request structure
	type PrescriptionRequest struct {
		PatientID            uuid.UUID `json:"patient_id"`
		Diagnosis            string    `json:"diagnosis"`
		PrescribedMedicineIDs []uint   `json:"prescribed_medicine_ids"`
		Status               string    `json:"status"`
	}

	// Parse the request body
	var req PrescriptionRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println(err.Error())
		return nil, fmt.Errorf("invalid request data")
	}

	// Validate required fields
	if req.PatientID == uuid.Nil {
		return nil, fmt.Errorf("patient_id is required")
	}
	if req.Diagnosis == "" {
		return nil, fmt.Errorf("diagnosis is required")
	}
	if len(req.PrescribedMedicineIDs) == 0 {
		return nil, fmt.Errorf("at least one prescribed_medicine_id is required")
	}
	// Fetch the medicines from the database
	var medicines []Medicine
	if err := db.Where("id IN ?", req.PrescribedMedicineIDs).Find(&medicines).Error; err != nil {
		return nil, fmt.Errorf("error fetching medicines")
	}

	// Check if all medicines were found
	if len(medicines) != len(req.PrescribedMedicineIDs) {
		return nil, fmt.Errorf("some medicines not found")
	}

	// Check if all medicines are in stock
	for _, medicine := range medicines {
		if !medicine.InStock {
			return nil, fmt.Errorf("medicine %s is out of stock", medicine.Name)
		}
	}

	// Create the prescription
	prescription := Prescription{
		ID:                 uuid.New(),
		PatientID:          req.PatientID,
		Diagnosis:          req.Diagnosis,
		PrescribedMedicines: medicines,
		Status:             req.Status,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Save the prescription to the database
	if err := db.Create(&prescription).Error; err != nil {
		return nil, fmt.Errorf("could not save prescription")
	}

	// Return the created prescription and nil error
	return &prescription, nil
}

// GetPrescriptions retrieves all prescriptions
func GetPrescriptions(c *fiber.Ctx)(*[]Prescription,error) {
	var prescriptions []Prescription

	if err := db.Preload("PrescribedMedicines").Find(&prescriptions).Error; err != nil {
		return nil,c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve prescriptions"})
	}
	return &prescriptions,nil
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
	if err := db.Preload("PrescribedMedicines").First(&prescription, "id = ?", id).Error; err != nil {
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

	// Handle PrescribedMedicines if provided in the request (optional, based on requirements)
	if len(updateData.PrescribedMedicines) > 0 {
		// Here, you might want to handle the medicines in the request:
		// - Add new medicines
		// - Remove old ones (or reset the list entirely, depending on your use case)
		// For example, you might just replace the existing medicines with the new ones provided.

		// This example just replaces the existing medicines:
		prescription.PrescribedMedicines = updateData.PrescribedMedicines
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
