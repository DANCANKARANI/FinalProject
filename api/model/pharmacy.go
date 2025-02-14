package model

import (
	"errors"
	"log"
	"strconv"
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

/*
edit inventory
@params inventoryID 
*/
func EditInventory(c *fiber.Ctx, inventoryID string) (*Inventory, error) {
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

	// Validate inventoryID
	if inventoryID == "" {
		log.Println("Inventory ID is missing")
		return nil, errors.New("inventory ID is required")
	}

	// Parse request body
	var req InventoryRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("Error parsing JSON:", err)
		return nil, errors.New("error parsing JSON data")
	}

	// Convert inventoryID to UUID
	inventoryUUID, err := uuid.Parse(inventoryID)
	if err != nil {
		log.Println("Invalid UUID format:", err)
		return nil, errors.New("invalid inventory ID format")
	}

	// Find the existing inventory item
	var inventory Inventory
	if err := db.First(&inventory, "id = ?", inventoryUUID).Error; err != nil {
		log.Println("Inventory item not found:", err)
		return nil, errors.New("inventory item not found")
	}

	// Update only provided fields
	if req.Name != "" {
		inventory.Name = req.Name
	}
	if req.Quantity > 0 {
		inventory.Quantity = req.Quantity
	}
	if req.Price > 0 {
		inventory.Price = req.Price
	}
	if req.Category != "" {
		inventory.Category = req.Category
	}
	if !req.ExpiryDate.IsZero() {
		inventory.ExpiryDate = req.ExpiryDate
	}
	if req.Supplier != "" {
		inventory.Supplier = req.Supplier
	}
	if req.Description != "" {
		inventory.Description = req.Description
	}
	if req.BatchNumber != "" {
		inventory.BatchNumber = req.BatchNumber
	}
	if req.ReorderLevel > 0 {
		inventory.ReorderLevel = req.ReorderLevel
	}

	// Set updated time
	inventory.UpdatedAt = time.Now()

	// Save the updated inventory item
	if err := db.Save(&inventory).Error; err != nil {
		log.Println("Error saving updated inventory:", err)
		return nil, errors.New("could not update inventory")
	}

	// Log and return updated inventory
	log.Println("Updated inventory:", inventory)
	return &inventory, nil
}

/*
delete inventory handler
@params inventoryID
*/
func DeleteInventory(c *fiber.Ctx, inventoryID string) error {
	// Get inventory ID from URL params
	
	if inventoryID == "" {
		return errors.New("inventory ID is required")
	}

	// Find the inventory item
	var inventory Inventory
	if err := db.First(&inventory, "id = ?", inventoryID).Error; err != nil {
		return errors.New("inventory item not found")
	}

	// Delete the inventory item
	if err := db.Delete(&inventory).Error; err != nil {
		return errors.New("could not delete inventory item")
	}

	return nil
}
//gets all the inventories

func GetInventory(c *fiber.Ctx) (*[]Inventory, error) {
	var inventories []Inventory

	// Get query parameters for filtering and pagination
	name := c.Query("name")         // Filter by name
	category := c.Query("category") // Filter by category
	limit, _ := strconv.Atoi(c.Query("limit", "10")) // Default limit: 10
	page, _ := strconv.Atoi(c.Query("page", "1"))    // Default page: 1
	offset := (page - 1) * limit

	// Build query with filters
	query := db.Model(&Inventory{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Fetch data with pagination
	if err := query.Limit(limit).Offset(offset).Find(&inventories).Error; err != nil {
		return nil, errors.New("could not retrieve inventory records")
	}

	return &inventories, nil
}
