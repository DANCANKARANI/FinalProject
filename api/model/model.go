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
	ID                 uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	PatientID          uuid.UUID      `json:"patient_id" gorm:"type:varchar(36);not null"` // Foreign key to Patient
	Patient            Patient        `json:"patient" gorm:"foreignKey:PatientID"`         // Relationship to Patient
	Diagnosis          string         `json:"diagnosis" gorm:"type:text"`
	PrescribedMedicines []Medicine     `json:"prescribed_medicines" gorm:"many2many:prescription_medicines"`
	Status             string         `json:"status" gorm:"type:varchar(20);default:'Pending'"`
	CreatedAt          time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

//reception struct represents a reception record
type Receptionist struct {
	gorm.Model
	FullName string `gorm:"not null"`
	Username string `gorm:"unique;not null"` // Add Username field
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}
type Patient struct {
	ID               uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	FirstName        string         `json:"first_name" gorm:"type:varchar(50);not null"`
	LastName         string         `json:"last_name" gorm:"type:varchar(50);not null"`
	Gender           string         `json:"gender" gorm:"type:varchar(10);not null"`
	DateOfBirth      time.Time      `json:"dob" gorm:"not null"`
	PatientNumber    string         `json:"patient_number" gorm:"type:varchar(20);unique;not null"`
	PhoneNumber      string         `json:"phone_number" gorm:"type:varchar(15);not null"`
	Email            string         `json:"email" gorm:"type:varchar(100)"`
	Address          string         `json:"address" gorm:"type:varchar(255)"`
	BloodGroup       string         `json:"blood_group" gorm:"type:varchar(10)"`
	MedicalHistory   string         `json:"medical_history" gorm:"type:text"`

	// Emergency-specific fields
	IsEmergency      bool           `json:"is_emergency" gorm:"default:false"` // Flag for emergency cases
	EmergencyContact string         `json:"emergency_contact" gorm:"type:varchar(15)"` // Contact person for emergencies
	TriageLevel      string         `json:"triage_level" gorm:"type:varchar(20)"` // Red, Yellow, Green
	InitialVitals    string         `json:"initial_vitals" gorm:"type:text"` // JSON or stringified vitals
	EmergencyNotes   string         `json:"emergency_notes" gorm:"type:text"` // Additional notes for emergency cases

	// Timestamps
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"` // Soft delete
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



// Referral represents a patient referral record
type Referral struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PatientID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"patient_id"`
	DoctorID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"doctor_id"`
	ReferredTo string         `gorm:"type:varchar(255);not null" json:"referred_to"`
	Reason     string         `gorm:"type:text;not null" json:"reason"`
	Diagnosis  string         `gorm:"type:text" json:"diagnosis"`
	LabResults string         `gorm:"type:jsonb" json:"lab_results"` // Store structured lab results in JSON format
	Status     string         `gorm:"type:enum('Pending', 'Accepted', 'Completed');default:'Pending'" json:"status"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
