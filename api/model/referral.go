package model

import (
	"fmt"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ReferPatient(c *fiber.Ctx,patient_id uuid.UUID) (*Referral, error) {
	// Define the request structure
	type ReferPatientRequest struct {
		PatientID   uuid.UUID `json:"patient_id"`   // ID of the patient being referred
		DoctorID    uuid.UUID `json:"doctor_id"`    // ID of the referring doctor
		ReferredTo  string    `json:"referred_to"`  // Name or ID of the receiving entity (e.g., specialist, hospital)
		Reason      string    `json:"reason"`       // Reason for the referral
		Diagnosis   string    `json:"diagnosis"`    // Diagnosis (optional)
		LabResults  string    `json:"lab_results"`  // Lab results in JSON format (optional)
	}

	// Parse the request body
	var req ReferPatientRequest
	if err := c.BodyParser(&req); err != nil {
		return nil, fmt.Errorf("invalid request data")
	}

	// Validate required fields
	id, err:= GetAuthUserID(c)
	if err != nil{
		fmt.Println(err.Error())
	}
	req.DoctorID=id
	fmt.Println(req.DoctorID)
	if  req.ReferredTo == "" || req.Reason == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	// Create the referral
	referral := Referral{
		ID:         uuid.New(),
		PatientID:  patient_id,
		DoctorID:   req.DoctorID,
		ReferredTo: req.ReferredTo,
		Reason:     req.Reason,
		Diagnosis:  req.Diagnosis,
		LabResults: req.LabResults,
		Status:     "Pending", // Default status
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save the referral to the database
	if err := db.Create(&referral).Error; err != nil {
		return nil, fmt.Errorf("failed to create referral")
	}

	// Return the created referral and nil error
	return &referral, nil
}