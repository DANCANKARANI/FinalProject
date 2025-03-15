package model

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)



/*
create inventory
*/
func CreateInventory(c *fiber.Ctx) (*Inventory, error) {
	// Define a struct for receiving request payload
	type InventoryRequest struct {
		Name         string `json:"name"`
		MedicineID   uint	`json:"medicine_id"`
		Quantity     int    `json:"quantity"`
		Category     string `json:"category"`
		ExpiryDate   string `json:"expiry_date"` // Use string for parsing
		ReorderLevel int    `json:"reorder_level"`
	}

	// Parse request body
	var req InventoryRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println(err.Error())
		return nil, errors.New("error parsing json data")
	}

	// Validate expiry_date
	if req.ExpiryDate == "" {
		return nil, errors.New("expiry date is required")
	}

	// Parse the expiry date
	expiryDate, err := time.Parse(time.RFC3339, req.ExpiryDate)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("invalid expiry date format: expected YYYY-MM-DDTHH:mm:ssZ")
	}

	// Create new inventory item
	inventory := Inventory{
		ID:           uuid.New(),
		Name:         req.Name,
		MedicineID: req.MedicineID,
		Quantity:     req.Quantity,
		Category:     req.Category,
		ExpiryDate:   expiryDate,
		ReorderLevel: req.ReorderLevel,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to database
	if err := db.Create(&inventory).Error; err != nil {
		return nil, errors.New("could not save inventory")
	}

	// Return success response
	return &inventory, nil
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
	if req.Category != "" {
		inventory.Category = req.Category
	}
	if !req.ExpiryDate.IsZero() {
		inventory.ExpiryDate = req.ExpiryDate
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
