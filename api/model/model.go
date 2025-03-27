package model

import (
	"time"
	"gorm.io/datatypes"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Medicine struct {
    ID                  uint          `json:"id" gorm:"primaryKey;autoIncrement"`
    Name                string        `json:"name" gorm:"type:varchar(100);not null"`
	Form                string         `json:"form" gorm:"type:varchar(50)"` // e.g., tablet, liquid, injection
    InStock             bool          `json:"in_stock" gorm:"default:true"`
    Inventories         []Inventory   `json:"inventories" gorm:"foreignKey:MedicineID"` // One-to-many relationship
    CreatedAt           time.Time     `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt           time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    		gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Inventory struct {
    ID           uuid.UUID      `json:"id" gorm:"type:varchar(36);primary_key"`
    MedicineID   uint           `json:"medicine_id" gorm:"not null"` // Foreign key to Medicine
    Name         string         `json:"name" gorm:"type:varchar(100);not null"`
    Quantity     int            `json:"quantity" gorm:"not null"`
    Category     string         `json:"category" gorm:"type:varchar(50)"`
    ExpiryDate   time.Time      `json:"expiry_date"`
    ReorderLevel int            `json:"reorder_level" gorm:"not null"`
    CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Prescription struct represents a prescription record.
type Prescription struct {
    ID                 uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	MedicineName			string		`json:"medicine_name" gorm:"type:varchar(100)"`
    PatientID          uuid.UUID      `json:"patient_id" gorm:"type:varchar(36);not null"` // Foreign key to Patient
    Patient            Patient        `json:"patient" gorm:"foreignKey:PatientID"`         // Relationship to Patient
    DoctorID           uuid.UUID      `json:"doctor_id" gorm:"type:varchar(36);not null"`   // Foreign key to Doctor
    Doctor             User           `json:"doctor" gorm:"foreignKey:DoctorID"`           // Relationship to Doctor
    Diagnosis          string         `json:"diagnosis" gorm:"type:text"`
	Dosage				string		`json:"dosage" gorm:"type:varchar(255)"`
	Instructions		string		`json:"instructions" gorm:"type:text"`
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
	Notes           []Note         `json:"notes" gorm:"foreignKey:PatientID"` // One-to-Many relationship

	// Timestamps
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"` // Soft delete
}

type LabTest struct {
	ID          uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	TestName    string         `json:"test_name" gorm:"type:varchar(100);not null"`
	Description string         `json:"description" gorm:"type:text"`
	SampleType  string         `json:"sample_type" gorm:"type:varchar(50);not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	PatientID   uuid.UUID      `json:"patient_id" gorm:"type:varchar(36);not null;index"`
	Patient     Patient        `json:"patient" gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Storing results as JSON
	Results datatypes.JSON `json:"results" gorm:"type:json"` // JSON field for test results

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}


// Referral represents a patient referral record
type Referral struct {
	ID          uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey;"`
	PatientID   uuid.UUID      `json:"patient_id" gorm:"type:varchar(36);"`
	DoctorID    uuid.UUID      `json:"doctor_id" gorm:"type:varchar(36);"`
	ReferredTo  string         `json:"referred_to" gorm:"type:varchar(255)"`
	Reason      string         `json:"reason" gorm:"type:text;"`
	Diagnosis   string         `json:"diagnosis" gorm:"type:text"`
	LabResults  string         `json:"lab_results" gorm:"type:text"` // Store structured lab results in JSON format
	Status      string         `json:"status" gorm:"type:varchar(20);default:'Pending';check:status IN ('Pending', 'Accepted', 'Completed')"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}


type Payment struct {
	ID              uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey;"` // Unique identifier for the payment
	BillingID       uuid.UUID `json:"billing_id" gorm:"type:varchar(36)"` // Foreign key to the Billing table
	PatientID       uuid.UUID `json:"patient_id" gorm:"type:varchar(36)"` // Foreign key to the Patients table
	Cost            float64   `json:"cost" gorm:"type:decimal(10,2);"`        // Cost of the payment
	PaymentMethod   string    `json:"payment_method" gorm:"type:varchar(50);"` // Payment method (e.g., M-Pesa, Credit Card)
	TransactionID   string    `json:"transaction_id" gorm:"type:varchar(100);"` // Transaction ID from the payment gateway
	PaymentStatus   string    `json:"payment_status" gorm:"type:varchar(50);"` // Payment status (e.g., Pending, Completed, Failed)
	CallbackURL     string    `json:"callback_url" gorm:"type:varchar(255);"`  // Callback URL for payment notifications
	CustomerPhone   string    `json:"customer_phone" gorm:"type:varchar(20);"` // Customer's phone number
	CustomerName    string    `json:"customer_name" gorm:"type:varchar(100);"` // Customer's name
	AccountReference string   `json:"account_reference" gorm:"type:varchar(100);"` // Account reference (e.g., order ID)
	TransactionDesc string    `json:"transaction_desc" gorm:"type:varchar(255);"` // Transaction description	
	TransactionDate string	  `json:"transaction_date" gorm:"type:varchar(255);"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;"` // Timestamp when the payment was created
	UpdatedAt       time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;"` // Timestamp when the payment was last updated
	// Relationships
	Billing Billing `json:"billing" gorm:"foreignKey:BillingID;references:ID"` // Many-to-one relationship with Billing
	Patient Patient `json:"patient" gorm:"foreignKey:PatientID;references:ID"` // Many-to-one relationship with Patient
}



type Billing struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey;"` // Unique identifier for the billing record
	PatientID   uuid.UUID `json:"patient_id" gorm:"type:varchar(36)"` // Foreign key to the Patients table
	Quantity    int       `json:"quantity"`           // Quantity of items or services billed
	Price       float64   `json:"price" gorm:"type:decimal(10,2);"` // Price per unit of the item or service
	Description string    `json:"description" gorm:"type:text;"` // Description of the bill (e.g., item/service details)
	Paid		bool	`json:"paid"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;"`     // Timestamp when the billing record was created
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;"`     // Timestamp when the billing record was last updated
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Patient  Patient   `json:"patient" gorm:"foreignKey:PatientID;references:ID"` // Many-to-one relationship with Patient
	Payments []Payment `json:"payments" gorm:"foreignKey:BillingID"` // One-to-many relationship with Payments
}

type ClinicBooking struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey;"` // Unique identifier for the billing record
    PatientID string         `gorm:"type:uuid;not null" json:"patient_id"`
    Reasons   string         `gorm:"type:text;not null" json:"reasons"`
    Date      time.Time      `gorm:"not null" json:"date"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    // Relationships
    Patient Patient `gorm:"foreignKey:PatientID;constraint:OnDelete:CASCADE" json:"patient"`
}

type Note struct {
	ID        uuid.UUID      `json:"id" gorm:"type:varchar(36);primaryKey"`
	Notes     string         `json:"notes" gorm:"type:text"`
	PatientID uuid.UUID      `json:"patient_id" gorm:"type:varchar(36);index"` // Foreign key
	Patient   Patient        `json:"-" gorm:"foreignKey:PatientID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"` // Soft delete
}