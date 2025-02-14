package model

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Inventory struct {
	ID           uuid.UUID      `json:"id" gorm:"type:varchar(36);primary_key"`
	Name         string         `json:"name" gorm:"type:varchar(100);not null"`
	Quantity     int            `json:"quantity" gorm:"not null"`
	Price        float64        `json:"price" gorm:"not null"`
	Category     string         `json:"category" gorm:"type:varchar(50)"`
	ExpiryDate   time.Time      `json:"expiry_date"`
	Supplier     string         `json:"supplier" gorm:"type:varchar(100)"`
	Description  string         `json:"description" gorm:"type:text"`
	BatchNumber  string         `json:"batch_number" gorm:"type:varchar(50);unique"`
	ReorderLevel int            `json:"reorder_level" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

/*
create inventory
*/
func CreateInventory(c *fiber.Ctx)(*Inventory, error) {
	// Define a struct for receiving request payload
	type InventoryRequest struct {
		Name         string    `json:"name"`
		Quantity     int       `json:"quantity"`
		Price        float64   `json:"price"`
		Category     string    `json:"category"`
		ExpiryDate   time.Time `json:"expiry_date"`
		Supplier     string    `json:"supplier"`
		Description  string    `json:"description"`
		BatchNumber  string    `json:"batch_number"`
		ReorderLevel int       `json:"reorder_level"`
	}

	// Parse request body
	var req InventoryRequest
	if err := c.BodyParser(&req); err != nil {
		return nil,errors.New("error parsing json data")
	}

	// Validate required fields
	if req.Name == "" || req.Quantity <= 0 || req.Price <= 0 {
		return nil,errors.New("name, quantity, and price are required")
	}

	// Create new inventory item
	inventory := Inventory{
		ID:           uuid.New(),
		Name:         req.Name,
		Quantity:     req.Quantity,
		Price:        req.Price,
		Category:     req.Category,
		ExpiryDate:   req.ExpiryDate,
		Supplier:     req.Supplier,
		Description:  req.Description,
		BatchNumber:  req.BatchNumber,
		ReorderLevel: req.ReorderLevel,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to database
	if err := db.Create(&inventory).Error; err != nil {
		return nil,errors.New("could not save inventory")
	}

	// Return success response
	return &inventory,nil
	
}
