package model

import (
	"errors"
	"log"
	"time"

	"github.com/dancankarani/medicare/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var db = database.ConnectDB()

// User struct
type User struct {
	ID          uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	FullName    string         `json:"full_name" gorm:"type:varchar(100)"`
	Email       string         `json:"email" gorm:"type:varchar(100);unique"`
	Username    string         `json:"username" gorm:"type:varchar(50);not null;unique"`
	PhoneNumber string         `json:"phone_number" gorm:"type:varchar(15)"`
	Role        string         `json:"role" gorm:"type:varchar(20)"` // Role: doctor, pharmacist, patient, etc.
	DateOfBirth time.Time      `json:"date_of_birth"`
	Password    string         `json:"password" gorm:"type:varchar(255)"`
	Address     string         `json:"address" gorm:"type:text"`
	Gender      string         `json:"gender" gorm:"type:varchar(10)"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func CreateUser(c *fiber.Ctx, user User) (*User, error) {
	// Generate a new UUID for the user
	user.ID = uuid.New()

	// Validate required fields
	if user.FullName == "" {
		return nil, errors.New("full name field should not be empty")
	}
	if user.Email == "" {
		return nil, errors.New("email field should not be empty")
	}
	if user.Username == "" {
		return nil, errors.New("username field should not be empty")
	}
	if user.Password == "" {
		return nil, errors.New("password field should not be empty")
	}
	if user.Role == "" {
		return nil, errors.New("role field should not be empty")
	}
	
	// Create the user in the database
	if err := db.Create(&user).Error; err != nil {
		log.Println("failed to create user:", err.Error())
		return nil, errors.New("failed to create user")
	}

	return &user, nil
}


//edit user
func EditUser(c *fiber.Ctx, userID uuid.UUID, updatedUser User) (*User, error) {
	var user User

	// Fetch the existing user
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		log.Println("user not found:", err.Error())
		return nil, errors.New("user not found")
	}

	// Validate required fields
	if updatedUser.FullName == "" {
		return nil, errors.New("full name field should not be empty")
	}
	if updatedUser.Email == "" {
		return nil, errors.New("email field should not be empty")
	}
	if updatedUser.Username == "" {
		return nil, errors.New("username field should not be empty")
	}
	if updatedUser.Role == "" {
		return nil, errors.New("role field should not be empty")
	}

	// Update user details
	user.FullName = updatedUser.FullName
	user.Email = updatedUser.Email
	user.Username = updatedUser.Username
	user.Role = updatedUser.Role
	if updatedUser.Password != "" { // Only update password if provided
		user.Password = updatedUser.Password
	}

	// Save the updated user
	if err := db.Save(&user).Error; err != nil {
		log.Println("failed to update user:", err.Error())
		return nil, errors.New("failed to update user")
	}

	return &user, nil
}

//get all users
func GetAllUsers(c *fiber.Ctx) ([]User, error) {
	var users []User

	// Fetch all users from the database
	if err := db.Find(&users).Error; err != nil {
		log.Println("failed to fetch users:", err.Error())
		return nil, errors.New("failed to fetch users")
	}

	return users, nil
}

//get user by role
func GetUsersByRole(c *fiber.Ctx, role string) ([]User, error) {
	var users []User

	// Fetch users by role
	if err := db.Where("role = ?", role).Find(&users).Error; err != nil {
		log.Println("failed to fetch users by role:", err.Error())
		return nil, errors.New("failed to fetch users by role")
	}

	return users, nil
}

//delete user
func DeleteUser(c *fiber.Ctx,userID string) error {
	// Get user ID from URL params
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Check if the user exists
	var user User
	if err := db.Where("id = ?", parsedID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Delete the user
	if err := db.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

//get one user by id
func GetOneUser(c *fiber.Ctx, userID string) (*User, error) {
	var user User

	// Fetch user from database
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		log.Println("error getting user:", err.Error())
		return nil, errors.New("user not found")
	}

	return &user, nil
}

