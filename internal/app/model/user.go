package model

import (
	"errors"
	"log"
	"time"

	"github.com/dancankarani/medicare/internal/app/database"
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
	UserID     		string         	`json:"user_id" gorm:"type:varchar(36);not null"`
	User       		User           	`gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Allergies		string			`json:"allergies" gorm:"type:text"`
	MedicalNotes	string			`json:"medical_notes" gorm:"type:text"`
	CreatedAt		time.Time		`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt		time.Time		`json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt		gorm.DeletedAt	`json:"deleted_at" gorm:"index"`
}

//doctor db struct model
type Doctor struct{
	ID				uuid.UUID 		`json:"id" gorm:"type:varchar(36)"`
	UserID     		string         	`json:"user_id" gorm:"type:varchar(36);not null"`
	User       		User           	`gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Speciality		string			`json:"speciality" gorm:"type:varchar(100)"`
	CreatedAt		time.Time		`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt		time.Time		`json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt		gorm.DeletedAt	`json:"deleted_at" gorm:"index"`
}

func CreateUser(c *fiber.Ctx,user User)error{
	user.ID=uuid.New()

	if (user.FullName  == "" ) {
		return errors.New("full name field should not be empty")
	}
	if err := db.Create(&user).Error; err != nil{
		log.Println("failed to create user:",err.Error())
		return errors.New("failed to create user")
	}
	return nil
}

//check if user already exist
