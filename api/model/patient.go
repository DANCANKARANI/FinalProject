package model

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreatePatient(c *fiber.Ctx) error {
	var patient Patient

	// Parse the request body into a map to handle custom fields
	var requestBody map[string]interface{}
	if err := c.BodyParser(&requestBody); err != nil {
		log.Println("error registering patient:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Manually parse `dob` based on its format
	if dobStr, ok := requestBody["dob"].(string); ok {
		var parsedDOB time.Time
		var err error

		// Try parsing as "YYYY-MM-DD"
		parsedDOB, err = time.Parse("2006-01-02", dobStr)
		if err != nil {
			// If that fails, try parsing as "YYYY-MM-DDTHH:MM:SSZ" (ISO 8601)
			parsedDOB, err = time.Parse(time.RFC3339, dobStr)
			if err != nil {
				log.Println("Error parsing date:", err)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Invalid date format",
					"message": "Expected date in YYYY-MM-DD or YYYY-MM-DDTHH:MM:SSZ format",
				})
			}
		}
		patient.DateOfBirth = parsedDOB
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing or invalid date of birth",
			"message": "Date of birth (dob) is required",
		})
	}

	// Parse the rest of the fields into the Patient struct
	if err := c.BodyParser(&patient); err != nil {
		log.Println("error registering patient:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validate phone number
	phone_number, err := utilities.ValidatePhoneNumber(patient.PhoneNumber, "KE")
	if err != nil {
		log.Println(err.Error())
		return utilities.ShowError(c, "Invalid phone number: "+phone_number, fiber.StatusBadRequest, nil)
	}
	patient.PhoneNumber = phone_number

	// Generate patient number
	patient.PatientNumber, err = generatePatientNumber(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate patient number"})
	}

	// Save to DB
	patient.ID = uuid.New()

	if err := db.Create(&patient).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create patient"})
	}

	billing := Billing{
		ID: uuid.New(),
		PatientID: patient.ID,
		Quantity: 1,
		Price: 20,
		Paid: false,
		Description: "Registration",
	}
	//create a billing table
	if err := db.Create(&billing).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create billing table"})
	}
	return c.Status(fiber.StatusCreated).JSON(patient)
}

//generate patient number

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
func UpdatePatient(c *fiber.Ctx) (*Patient,error) {
	patientID := c.Params("id")

	// Find the existing patient by ID
	var patient Patient
	if err := db.First(&patient, "id = ?", patientID).Error; err != nil {
		return nil,errors.New("patient not found")
	}

	// Parse request body
	updateData := new(Patient)
	if err := c.BodyParser(updateData); err != nil {
		return nil,errors.New("failed to parse json data")
	}

	// Initialize error map
	error := make(map[string][]string)

	// Validate Phone Number
	if updateData.PhoneNumber != "" {
		phone_number, err := utilities.ValidatePhoneNumber(updateData.PhoneNumber, "KE")
		if err != nil {
			error["phone_number"] = append(error["phone_number"], err.Error())
		} else {
			patient.PhoneNumber = phone_number // Update phone number
		}
	}

	// Validate Full Name (Ensure it's not empty)
	if updateData.FirstName == "" || updateData.LastName == "" {
		error["full_name"] = append(error["full_name"], "Full name is required")
	} else {
		patient.FirstName = updateData.FirstName
		patient.LastName = updateData.LastName
	}

	// Validate Date of Birth (Should be a valid date and in the past)
	if !updateData.DateOfBirth.IsZero() {
		if updateData.DateOfBirth.After(time.Now()) {
			error["date_of_birth"] = append(error["date_of_birth"], "Date of birth must be in the past")
		} else {
			patient.DateOfBirth = updateData.DateOfBirth
		}
	}

	// Validate Medical Notes (Optional but should not exceed 500 characters)
	if len(updateData.MedicalHistory) > 500 {
		error["medical_notes"] = append(error["medical_notes"], "Medical notes should not exceed 500 characters")
	} else {
		patient.MedicalHistory = updateData.MedicalHistory
	}

	// If there are validation errors, return them
	if len(error) > 0 {
		return nil,errors.New("validation failed")
	}

	// Update the patient's record
	if err := db.Save(&patient).Error; err != nil {
		return nil, errors.New("failed to update patient")
	}

	// Return updated patient details
	return &patient,nil
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
	limit, _ := strconv.Atoi(c.Query("limit", "100"))  // Default limit = 10
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
func GetPatientBills(patientID uuid.UUID) (*[]Billing, error) {
    var bills []Billing
    err := db.Where("patient_id = ?", patientID).Find(&bills).Error
    if err != nil {
        return nil, err
    }
    return &bills, nil
}

func BookClinic(c *fiber.Ctx) error {
	// Parse request body
	var request struct {
		PatientID string    `json:"patient_id"`
		Reasons   string    `json:"reasons"`
		Date      time.Time `json:"date"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate input fields
	if  request.Reasons == "" || request.Date.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All fields are required"})
	}

	// Check if patient exists
	request.PatientID= c.Params("id")
	var patient Patient
	if err := db.First(&patient, "id = ?", request.PatientID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Patient not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Create a new clinic booking
	booking := ClinicBooking{
		ID:        uuid.New(),
		PatientID: request.PatientID,
		Reasons:   request.Reasons,
		Date:      request.Date,
	}

	// Save the booking to the database
	if err := db.Create(&booking).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create booking"})
	}

	// Success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Clinic appointment booked successfully",
		"booking": booking,
	})
}

func GetAllBillings() (*[]Billing, error) {
    var bills []Billing
    err := db.Find(&bills).Error
    if err != nil {
        return nil, err
    }
    return &bills, nil
}
