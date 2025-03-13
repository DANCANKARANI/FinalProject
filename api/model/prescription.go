package model

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreatePrescription handles creating a new prescription
func CreatePrescription(c *fiber.Ctx)(*Prescription,error){
	type PrescriptionRequest struct {
		PatientName         string   `json:"patient_name"`
		Age                 int      `json:"age"`
		Diagnosis           string   `json:"diagnosis"`
		PrescribedMedicineIDs []uint `json:"prescribed_medicine_ids"`
		Status              string   `json:"status"`
	}

	var req PrescriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return nil,c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data"})
	}

	if req.PatientName == "" || req.Diagnosis == "" || len(req.PrescribedMedicineIDs) == 0 {
		return nil,c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Patient name, diagnosis, and medicines are required"})
	}

	var medicines []Medicine
	if err := db.Where("id IN ?", req.PrescribedMedicineIDs).Find(&medicines).Error; err != nil {
		return nil,c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching medicines"})
	}

	if len(medicines) != len(req.PrescribedMedicineIDs) {
		return nil,c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Some medicines not found"})
	}

	for _, medicine := range medicines {
		if !medicine.InStock {
			return nil,c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Some medicines are out of stock"})
		}
	}

	prescription := Prescription{
		ID:                 uuid.New(),
		PatientName:        req.PatientName,
		Age:                req.Age,
		Diagnosis:          req.Diagnosis,
		PrescribedMedicines: medicines,
		Status:             req.Status,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := db.Create(&prescription).Error; err != nil {
		return nil,c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save prescription"})
	}
	return &prescription,nil
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
	if updateData.PatientName != "" {
		prescription.PatientName = updateData.PatientName
	}
	if updateData.Age != 0 {
		prescription.Age = updateData.Age
	}
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
