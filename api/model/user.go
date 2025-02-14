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
//user db struct model
type User struct {
	ID          uuid.UUID		`json:"id" gorm:"type:varchar(36)"`
	FullName    string 			`json:"full_name" gorm:"type:varchar(100)"`
	Email       string 			`json:"email" gorm:"type:varchar(100)"`
	Username	string			`json:"username" gorm:"type:varchar(50);not Null"`
	PhoneNumber string 			`json:"phone_number" gorm:"type:varchar(15)"`
	DateOfBirth time.Time		`json:"date_of_birth"`
	Password	string			`json:"password" gorm:"type:varchar(255)"`
	Adrress		string			`json:"address" gorm:"type:text"`
	Gender		string			`json:"gender" gorm:"type:varchar(10)"`
	CreatedAt	time.Time		`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt	time.Time		`json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	gorm.DeletedAt	`json:"deleted_at" gorm:"index"`
}
//patients db struct model
type Patient struct{
	ID				uuid.UUID 		`json:"id" gorm:"type:varchar(36)"`
	FullName    	string 			`json:"full_name" gorm:"type:varchar(100)"`
	DateOfBirth 	time.Time		`json:"date_of_birth"`
	PatientNumber 	string			`json:"patient_number" gorm:"type:varchar(20)"`
	PhoneNumber 	string 			`json:"phone_number" gorm:"type:varchar(15)"`
	MedicalNotes	string			`json:"medical_notes" gorm:"type:text"`
	CreatedAt		time.Time		`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt		time.Time		`json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt		gorm.DeletedAt	`json:"deleted_at" gorm:"index"`
}

//pharmacy db struct model
type Pharmacist struct {
	ID        uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:varchar(36)"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

//doctor db struct model
type Doctor struct{
	ID				uuid.UUID 		`json:"id" gorm:"type:varchar(36)"`
	UserID     		uuid.UUID     	`json:"user_id" gorm:"type:varchar(36)"`
	User       		User           	`gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Speciality		string			`json:"speciality" gorm:"type:varchar(100)"`
	CreatedAt		time.Time		`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt		time.Time		`json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt		gorm.DeletedAt	`json:"deleted_at" gorm:"index"`
}

func CreateUser(c *fiber.Ctx, user User, isPharmacist bool, isDoctor bool, specialty string) (*User, error) {
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

	// Create the user in the database
	if err := db.Create(&user).Error; err != nil {
		log.Println("failed to create user:", err.Error())
		return nil, errors.New("failed to create user")
	}

	// Prepare response data
	responseData := fiber.Map{
		"user": user,
	}

	// If the user should be a pharmacist, create a pharmacist entry
	if isPharmacist {
		pharmacist := Pharmacist{
			ID:     uuid.New(),
			UserID: user.ID,
			User:   user, // Include the user details
		}

		if err := db.Create(&pharmacist).Error; err != nil {
			log.Println("failed to assign pharmacist role:", err.Error())
			return nil, errors.New("failed to assign pharmacist role")
		}

		responseData["pharmacist"] = pharmacist
	}

	// If the user should be a doctor, create a doctor entry
	if isDoctor {
		if specialty == "" {
			return nil, errors.New("specialty field should not be empty for a doctor")
		}

		doctor := Doctor{
			ID:         uuid.New(),
			UserID:     user.ID,
			Speciality: specialty,
			User:       user, // Include the user details
		}

		if err := db.Create(&doctor).Error; err != nil {
			log.Println("failed to assign doctor role:", err.Error())
			return nil, errors.New("failed to assign doctor role")
		}

		responseData["doctor"] = doctor
	}

	return &user, nil
}


func EditUser(c *fiber.Ctx, userID uuid.UUID, updatedUser User, isPharmacist bool, isDoctor bool, specialty string) (*User, error) {
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

    // Update user details
    user.FullName = updatedUser.FullName
    user.Email = updatedUser.Email
    user.Username = updatedUser.Username
    if updatedUser.Password != "" { // Only update password if provided
        user.Password = updatedUser.Password
    }

    if err := db.Save(&user).Error; err != nil {
        log.Println("failed to update user:", err.Error())
        return nil, errors.New("failed to update user")
    }

    // Prepare response data
    responseData := fiber.Map{
        "user": user,
    }

    // Handle pharmacist role
    if isPharmacist {
        var pharmacist Pharmacist
        if err := db.First(&pharmacist, "user_id = ?", user.ID).Error; err != nil {
            pharmacist = Pharmacist{ID: uuid.New(), UserID: user.ID}
            if err := db.Create(&pharmacist).Error; err != nil {
                log.Println("failed to assign pharmacist role:", err.Error())
                return nil, errors.New("failed to assign pharmacist role")
            }
        }
        responseData["pharmacist"] = pharmacist
    } else {
        db.Where("user_id = ?", user.ID).Delete(&Pharmacist{})
    }

    // Handle doctor role
    if isDoctor {
        if specialty == "" {
            return nil, errors.New("specialty field should not be empty for a doctor")
        }

        var doctor Doctor
        if err := db.First(&doctor, "user_id = ?", user.ID).Error; err != nil {
            doctor = Doctor{ID: uuid.New(), UserID: user.ID, Speciality: specialty}
            if err := db.Create(&doctor).Error; err != nil {
                log.Println("failed to assign doctor role:", err.Error())
                return nil, errors.New("failed to assign doctor role")
            }
        } else {
            doctor.Speciality = specialty
            if err := db.Save(&doctor).Error; err != nil {
                log.Println("failed to update doctor role:", err.Error())
                return nil, errors.New("failed to update doctor role")
            }
        }
        responseData["doctor"] = doctor
    } else {
        db.Where("user_id = ?", user.ID).Delete(&Doctor{})
    }

    return &user, nil
}

/*Get all the Pharmacists*/
func GetPharmacists(c *fiber.Ctx) ([]User, error) {
    var pharmacists []Pharmacist
    var users []User

    // Fetch all pharmacists
    if err := db.Find(&pharmacists).Error; err != nil {
        log.Println("failed to fetch pharmacists:", err.Error())
        return nil, errors.New("failed to fetch pharmacists")
    }

    // Get associated user data
    for _, pharmacist := range pharmacists {
        var user User
        if err := db.First(&user, "id = ?", pharmacist.UserID).Error; err == nil {
            users = append(users, user)
        }
    }

    return users, nil
}

/*Get all the doctors*/
func GetDoctors(c *fiber.Ctx) ([]User, error) {
    var doctors []Doctor
    var users []User

    // Fetch all doctors
    if err := db.Find(&doctors).Error; err != nil {
        log.Println("failed to fetch doctors:", err.Error())
        return nil, errors.New("failed to fetch doctors")
    }

    // Get associated user data
    for _, doctor := range doctors {
        var user User
        if err := db.First(&user, "id = ?", doctor.UserID).Error; err == nil {
            users = append(users, user)
        }
    }

    return users, nil
}

/*Get all users*/
func GetAllUsers(c *fiber.Ctx) ([]User, error) {
    var users []User

    // Fetch all users from the database
    if err := db.Find(&users).Error; err != nil {
        log.Println("failed to fetch users:", err.Error())
        return nil, errors.New("failed to fetch users")
    }

    return users, nil
}

/*Delete user*/
func DeleteUser(c *fiber.Ctx) error {
	// Get user ID from URL params
	userID := c.Params("id")
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

	// Delete related roles (Pharmacist and Doctor)
	db.Where("user_id = ?", parsedID).Delete(&Pharmacist{})
	db.Where("user_id = ?", parsedID).Delete(&Doctor{})

	// Delete the user
	if err := db.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User and associated roles deleted successfully",
	})
}

/*
get user by id
@params user_id string
*/
func GetOneUser(uc *fiber.Ctx, userID string) (*User, error) {
	var user User

	// Fetch user from database
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		log.Println("error geting user:"+err.Error())
		return nil, errors.New("errors getting user with id:"+userID)
	}

	return &user, nil
}