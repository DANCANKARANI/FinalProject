package model

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreatePatient(c *fiber.Ctx) error {
	var patient Patient
	if err := c.BodyParser(&patient); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	errors := make(map[string][]string)

	phone_number, err := utilities.ValidatePhoneNumber(patient.PhoneNumber, "KE")
	if err != nil {
		errors["phone_number"] = append(errors["phone_number"], err.Error())
		return utilities.ShowError(c, "Invalid phone number", fiber.StatusBadRequest, errors)
	}

	patient.ID = uuid.New()
	patient.PhoneNumber = phone_number
	patient.PatientNumber,err = generatePatientNumber(db) // Function to auto-generate unique patient numbers
	if err != nil{
		errors["patient_number"] = append(errors["patient_number"], err.Error())
	}

	if err := db.Create(&patient).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create patient"})
	}

	return c.Status(fiber.StatusCreated).JSON(patient)
}

func generatePatientNumber(db *gorm.DB) (string, error) {
	var lastPatient Patient
	if err := db.Order("created_at desc").First(&lastPatient).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	// Default start number if no patients exist
	nextNumber := 1

	// Extract last sequential number and increment it
	if lastPatient.PatientNumber != "" {
		var lastSeq int
		fmt.Sscanf(lastPatient.PatientNumber, "PAT-%d-%04d", new(int), &lastSeq)
		nextNumber = lastSeq + 1
	}

	// Generate the new patient number with format: "PAT-2025-0001"
	newPatientNumber := fmt.Sprintf("PAT-%d-%04d", time.Now().Year(), nextNumber)

	return newPatientNumber, nil
}



//update patients details
func UpdatePatient(c *fiber.Ctx) error {
	patientID := c.Params("id")

	// Find the existing patient by ID
	var patient Patient
	if err := db.First(&patient, "id = ?", patientID).Error; err != nil {
		return utilities.ShowError(c, "Patient not found", fiber.StatusNotFound, nil)
	}

	// Parse request body
	updateData := new(Patient)
	if err := c.BodyParser(updateData); err != nil {
		return utilities.ShowError(c, "Invalid request format", fiber.StatusBadRequest, nil)
	}

	// Initialize error map
	errors := make(map[string][]string)

	// Validate Phone Number
	if updateData.PhoneNumber != "" {
		phone_number, err := utilities.ValidatePhoneNumber(updateData.PhoneNumber, "KE")
		if err != nil {
			errors["phone_number"] = append(errors["phone_number"], err.Error())
		} else {
			patient.PhoneNumber = phone_number // Update phone number
		}
	}

	// Validate Full Name (Ensure it's not empty)
	if updateData.FullName == "" {
		errors["full_name"] = append(errors["full_name"], "Full name is required")
	} else {
		patient.FullName = updateData.FullName
	}

	// Validate Date of Birth (Should be a valid date and in the past)
	if !updateData.DateOfBirth.IsZero() {
		if updateData.DateOfBirth.After(time.Now()) {
			errors["date_of_birth"] = append(errors["date_of_birth"], "Date of birth must be in the past")
		} else {
			patient.DateOfBirth = updateData.DateOfBirth
		}
	}

	// Validate Medical Notes (Optional but should not exceed 500 characters)
	if len(updateData.MedicalNotes) > 500 {
		errors["medical_notes"] = append(errors["medical_notes"], "Medical notes should not exceed 500 characters")
	} else {
		patient.MedicalNotes = updateData.MedicalNotes
	}

	// If there are validation errors, return them
	if len(errors) > 0 {
		return utilities.ShowError(c, "Validation failed", fiber.StatusBadRequest, errors)
	}

	// Update the patient's record
	if err := db.Save(&patient).Error; err != nil {
		return utilities.ShowError(c, "Failed to update patient", fiber.StatusInternalServerError, nil)
	}

	// Return updated patient details
	return utilities.ShowSuccess(c,"patient updated successfully",1,patient)
}



/*delete patient*/
func DeletePatient(c *fiber.Ctx) error {
	patientID := c.Params("id")

	// Check if the provided ID is a valid UUID
	id, err := uuid.Parse(patientID)
	if err != nil {
		return utilities.ShowError(c, "Invalid patient ID", fiber.StatusBadRequest, nil)
	}

	// Find the patient
	var patient Patient
	if err := db.First(&patient, "id = ?", id).Error; err != nil {
		return utilities.ShowError(c, "Patient not found", fiber.StatusNotFound, nil)
	}

	// Delete the patient (soft delete)
	if err := db.Delete(&patient).Error; err != nil {
		return utilities.ShowError(c, "Failed to delete patient", fiber.StatusInternalServerError, nil)
	}

	// Success response
	return utilities.ShowMessage(c,"patient deleted",1)
}

/*get all Patients*/
func GetPatients(c *fiber.Ctx) error {
	// Pagination parameters
	limit, _ := strconv.Atoi(c.Query("limit", "10"))  // Default limit = 10
	page, _ := strconv.Atoi(c.Query("page", "1"))    // Default page = 1
	offset := (page - 1) * limit

	// Search filter (optional)
	search := c.Query("search")

	// Query patients from database
	var patients []Patient
	query := db.Model(&Patient{})

	// Apply search filter (FullName, PhoneNumber)
	if search != "" {
		query = query.Where("full_name LIKE ? OR phone_number LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Fetch records with pagination
	if err := query.Limit(limit).Offset(offset).Find(&patients).Error; err != nil {
		return utilities.ShowError(c, "Failed to fetch patients", fiber.StatusInternalServerError, nil)
	}

	// Count total patients
	var total int64
	db.Model(&Patient{}).Count(&total)

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Patients retrieved successfully",
		"data":    patients,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
