package model

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateNote(c *fiber.Ctx) error {
	note := new(Note)
	if err := c.BodyParser(note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	note.ID = uuid.New()

	if err := db.Create(&note).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create note"})
	}

	return c.Status(fiber.StatusCreated).JSON(note)
}

// Get all notes
func GetNotes(c *fiber.Ctx) error {
	var notes []Note
	db.Find(&notes)
	return c.JSON(notes)
}

// Get a single note by ID
func GetNote(c *fiber.Ctx) error {
	id := c.Params("id")
	var note Note

	if err := db.First(&note, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Note not found"})
	}

	return c.JSON(note)
}

// Update a note
func UpdateNote(c *fiber.Ctx) error {
	id := c.Params("id")
	var note Note

	if err := db.First(&note, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Note not found"})
	}

	if err := c.BodyParser(&note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := db.Save(&note).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update note"})
	}

	return c.JSON(note)
}