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
	UserID    uuid.UUID `json:"user_id" gorm:"type:varchar(36);not null;unique"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

//doctor db struct model
type Doctor struct{
	ID				uuid.UUID 		`json:"id" gorm:"type:varchar(36)"`
	UserID     		uuid.UUID     	`json:"user_id" gorm:"type:varchar(36);not null"`
	User       		User           	`gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Speciality		string			`json:"speciality" gorm:"type:varchar(100)"`
	CreatedAt		time.Time		`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt		time.Time		`json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt		gorm.DeletedAt	`json:"deleted_at" gorm:"index"`
}

func CreateUser(c *fiber.Ctx, user User, isPharmacist bool, isDoctor bool, specialty string) error {
	user.ID = uuid.New()

	// Validate required fields
	if user.FullName == "" {
		return errors.New("full name field should not be empty")
	}
	if user.Email == "" {
		return errors.New("email field should not be empty")
	}
	if user.Username == "" {
		return errors.New("username field should not be empty")
	}
	if user.Password == "" {
		return errors.New("password field should not be empty")
	}

	// Create the user in the database
	if err := db.Create(&user).Error; err != nil {
		log.Println("failed to create user:", err.Error())
		return errors.New("failed to create user")
	}

	// If the user should be a pharmacist, create a pharmacist entry
	if isPharmacist {
		pharmacist := Pharmacist{
			ID:     uuid.New(),
			UserID: user.ID,
		}

		if err := db.Create(&pharmacist).Error; err != nil {
			log.Println("failed to assign pharmacist role:", err.Error())
			return errors.New("failed to assign pharmacist role")
		}
	}

	// If the user should be a doctor, create a doctor entry
	if isDoctor {
		if specialty == "" {
			return errors.New("specialty field should not be empty for a doctor")
		}

		doctor := Doctor{
			ID:        uuid.New(),
			UserID:    user.ID,
			Speciality: specialty,
		}

		if err := db.Create(&doctor).Error; err != nil {
			log.Println("failed to assign doctor role:", err.Error())
			return errors.New("failed to assign doctor role")
		}
	}

	return nil
}


//check if user already exist
