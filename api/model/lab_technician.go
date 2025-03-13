package model

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func UploadLabTest(c *fiber.Ctx) (*LabTest, error) {
	// Parse the request body into a LabTest struct
	var labTest LabTest
	if err := c.BodyParser(&labTest); err != nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if labTest.TestName == "" || labTest.Cost <= 0 || labTest.PatientID == uuid.Nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Test name, cost, and patient ID are required",
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

	// Return the created lab test
	return &labTest, c.Status(fiber.StatusCreated).JSON(labTest)
}