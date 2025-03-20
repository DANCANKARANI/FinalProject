package model

import (
	"encoding/json"
	"errors"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func UploadLabTest(c *fiber.Ctx) (*LabTest, error) {
	// Parse the request body into a LabTest struct
	var labTest LabTest
	if err := c.BodyParser(&labTest); err != nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if the patient exists
	var patient Patient
	if err := db.First(&patient, "id = ?", labTest.PatientID).Error; err != nil {
		return nil, c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Patient not found",
		})
	}

	// Generate a new UUID for the lab test
	labTest.ID = uuid.New()

	// Set timestamps
	labTest.CreatedAt = time.Now()
	labTest.UpdatedAt = time.Now()

	// Save the lab test to the database
	if err := db.Create(&labTest).Error; err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save lab test to the database",
		})
	}
	
	billing := Billing{
		ID: uuid.New(),
		PatientID: patient.ID,
		Quantity: 1,
		Price: 100,
		Paid: false,
		Description: "Lab Test",
	}
	//create a billing table
	if err := db.Create(&billing).Error; err != nil {
		return nil,errors.New("failed to add bill")
	}

	// Return the created lab test
	return &labTest, c.Status(fiber.StatusCreated).JSON(labTest)
}



// GetAllLabTest retrieves all lab tests from the database
func GetAllLabTest(c *fiber.Ctx) (*[]LabTest, error) {
	var labTests []LabTest

	// Fetch all lab tests from the database and preload the Patient relationship
	if err := db.Preload("Patient").Find(&labTests).Error; err != nil {
		return nil, errors.New("failed to fetch data:"+err.Error()) // Return nil and the error
	}

	// Return the list of lab tests and nil error
	return &labTests, nil
}

// GetLabTestByID retrieves a specific lab test by its ID
func GetLabTestByID(c *fiber.Ctx, lab_test_id uuid.UUID) (*LabTest, error) {
	var labTest LabTest

	// Fetch the lab test by ID from the database
	if err := db.First(&labTest, "id = ?", lab_test_id).Error; err != nil {
		return nil, c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Lab test not found",
		})
	}

	// Return the lab test
	return &labTest, c.Status(fiber.StatusOK).JSON(labTest)
}


func CreateLabTestResult(c *fiber.Ctx) error {
	var resultRequest struct {
		LabTestID uuid.UUID `json:"lab_test_id"`
		Result    string    `json:"result"`
		Remarks   string    `json:"remarks"`
		TestedBy  uuid.UUID `json:"tested_by"`
	}

	// Parse the request body
	if err := c.BodyParser(&resultRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Get the authenticated user ID
	user_id, _ := GetAuthUserID(c)

	// Start a database transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Fetch the existing LabTest record
	var labTest LabTest
	if err := tx.First(&labTest, "id = ?", resultRequest.LabTestID).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Lab test not found"})
	}

	// Unmarshal the existing results (if any)
	var results map[string]interface{}
	if len(labTest.Results) > 0 {
		if err := json.Unmarshal(labTest.Results, &results); err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to parse existing results"})
		}
	} else {
		results = make(map[string]interface{})
	}

	// Add the new result to the results map
	results[time.Now().Format(time.RFC3339)] = map[string]interface{}{
		"result":    resultRequest.Result,
		"remarks":   resultRequest.Remarks,
		"tested_by": user_id,
	}

	// Marshal the updated results back to JSON
	updatedResultsJSON, err := json.Marshal(results)
	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to marshal updated results"})
	}

	// Update the LabTest record with the new results and set is_active to false
	labTest.Results = datatypes.JSON(updatedResultsJSON)
	labTest.IsActive = false // Set is_active to false

	// Save the updated LabTest record
	if err := tx.Save(&labTest).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update lab test results"})
	}

	// Commit the transaction
	tx.Commit()

	// Return the updated LabTest record
	return c.Status(200).JSON(labTest)
}