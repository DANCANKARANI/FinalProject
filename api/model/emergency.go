package model

import (
	"github.com/gofiber/fiber/v2"
)

func GetEmergencyCases(c *fiber.Ctx) (*[]Patient, error) {
	var emergencyPatients []Patient

	// Query the database for emergency cases
	if err := db.Where("is_emergency = ?", true).Find(&emergencyPatients).Error; err != nil {
		// Return nil and the error if the query fails
		return nil, err
	}

	// Return the list of emergency patients and nil error
	return &emergencyPatients, nil
}