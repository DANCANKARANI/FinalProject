package model

import (
	"time"
	"github.com/gofiber/fiber/v2"
)

// CreateMedicine adds a new medicine to the database
func CreateMedicine(c *fiber.Ctx) (*Medicine, error) {
	// Parse the request body into a Medicine struct
	var req Medicine
	if err := c.BodyParser(&req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate the request
	if req.Name == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Medicine name is required")
	}

	// Check if medicine already exists
	var existingMedicine Medicine
	if err := db.Where("name = ?", req.Name).First(&existingMedicine).Error; err == nil {
		return nil, fiber.NewError(fiber.StatusConflict, "Medicine already exists")
	}

	// Set timestamps
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	// Save to database
	if err := db.Create(&req).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to save medicine")
	}

	// Return the created medicine
	return &req, nil
}


// GetMedicines retrieves all medicines from the database
func GetMedicines(c *fiber.Ctx)(*[]Medicine, error) {
	var medicines []Medicine

	if err := db.Find(&medicines).Error; err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve medicines"})
	}

	return &medicines,nil
}


// GetMedicine retrieves a single medicine by ID
func GetMedicine(c *fiber.Ctx,id string) (*Medicine,error) {
	
	var medicine Medicine

	if err := db.First(&medicine, id).Error; err != nil {
		return nil,c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Medicine not found"})
	}

	return &medicine,nil
}


// UpdateMedicine updates an existing medicine
func UpdateMedicine(c *fiber.Ctx, id string) (*Medicine, error) {
	var medicine Medicine

	// Check if medicine exists
	if err := db.First(&medicine, id).Error; err != nil {
		return nil, c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Medicine not found"})
	}

	// Parse request body
	var updateData Medicine
	if err := c.BodyParser(&updateData); err != nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data"})
	}

	// Update only the fields that are provided in the request
	if updateData.Name != "" {
		medicine.Name = updateData.Name
	}

	
	if updateData.Form != "" {
		medicine.Form = updateData.Form
	}

	// Assuming "InStock" is a boolean and you can update it based on provided data
	medicine.InStock = updateData.InStock

	// Update the timestamp for the last update
	medicine.UpdatedAt = time.Now()

	// Save updated medicine
	if err := db.Save(&medicine).Error; err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update medicine"})
	}

	return &medicine, nil
}



// DeleteMedicine deletes a medicine by ID
func DeleteMedicine(c *fiber.Ctx,id string) error {
	var medicine Medicine

	// Check if medicine exists
	if err := db.First(&medicine, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Medicine not found"})
	}

	// Delete medicine
	if err := db.Delete(&medicine).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete medicine"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Medicine deleted successfully"})
}
