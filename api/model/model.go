package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Medicine struct represents a medicine with stock information.
type Medicine struct {
	ID                  uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name                string        `json:"name" gorm:"type:varchar(100);not null"`
	Dosage              string        `json:"dosage" gorm:"type:varchar(50)"`
	Frequency           string        `json:"frequency" gorm:"type:varchar(50)"`
	Duration            string        `json:"duration" gorm:"type:varchar(50)"`
	Route               string        `json:"route" gorm:"type:varchar(50)"`
	SpecialInstructions string        `json:"special_instructions" gorm:"type:text"`
	InStock             bool          `json:"in_stock" gorm:"default:true"`
	Prescriptions       []*Prescription `json:"-" gorm:"many2many:prescription_medicines"`
	CreatedAt           time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

// Prescription struct represents a prescription record.
type Prescription struct {
	ID                 uuid.UUID  `json:"id" gorm:"type:varchar(36);primaryKey"`
	PatientName        string     `json:"patient_name" gorm:"type:varchar(100);not null"`
	Age                int        `json:"age"`
	Diagnosis          string     `json:"diagnosis" gorm:"type:text"`
	PrescribedMedicines []Medicine `json:"prescribed_medicines" gorm:"many2many:prescription_medicines"`
	Status            string     `json:"status" gorm:"type:varchar(20);default:'Pending'"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
type Receptionist struct {
	gorm.Model
	FullName string `gorm:"not null"`
	Username string `gorm:"unique;not null"` // Add Username field
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}
type Patient struct {
	ID               uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	FirstName        string         `json:"first_name" gorm:"type:varchar(50)"`
	LastName         string         `json:"last_name" gorm:"type:varchar(50)"`
	Gender           string         `json:"gender" gorm:"type:varchar(10)"`
	DateOfBirth      time.Time      `json:"dob"`
	PatientNumber  string         `json:"patient_number" gorm:"type:varchar(20);unique"`
	PhoneNumber      string         `json:"phone_number" gorm:"type:varchar(15)"`
	Email            string         `json:"email" gorm:"type:varchar(100);unique"`
	Address          string         `json:"address" gorm:"type:varchar(255)"`
	EmergencyContact string         `json:"emergency_contact" gorm:"type:varchar(15)"`
	BloodGroup       string         `json:"blood_group" gorm:"type:varchar(10)"`
	MedicalHistory   string         `json:"medical_history" gorm:"type:text"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
type LabTest struct {
	ID          uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	TestName    string         `json:"test_name" gorm:"type:varchar(100);not null"`
	Description string         `json:"description" gorm:"type:text"`
	Cost        float64        `json:"cost" gorm:"type:decimal(10,2);not null"`
	Duration    string         `json:"duration" gorm:"type:varchar(50);not null"`
	SampleType  string         `json:"sample_type" gorm:"type:varchar(50);not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	PatientID   uuid.UUID      `json:"patient_id" gorm:"type:varchar(36);not null"` // Foreign key to Patient
	Patient     Patient        `json:"patient" gorm:"foreignKey:PatientID"`         // Relationship to Patient
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}