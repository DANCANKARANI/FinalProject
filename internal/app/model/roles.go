package model

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          uuid.UUID      `json:"id" gorm:"type:varchar(36)"`
	RoleName    string         `json:"role_name" gorm:"type:varchar(50)"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

//creates new role
func CreateRole(role Role)(*Role,error){
	db.AutoMigrate(&Role{})
	role.ID=uuid.New()
	if err := db.Create(&role).Error; err != nil {
		log.Println("error creating role:",err.Error())
		return nil,errors.New("error creating role")
	}
	return &role,nil
}

// UpdateRole updates an existing role in the database by ID
func UpdateRole(roleID uuid.UUID, role Role) (*Role, error) {
	// Find the role by ID
	
	if err := db.First(&role, "id = ?", roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("role not found:", err.Error())
			return nil, errors.New("role not found")
		}
		log.Println("error finding role:", err.Error())
		return nil, errors.New("error finding role")
	}

	// Update the fields that are allowed to be updated

	if err := db.Updates(&role).Error; err != nil {
		log.Println("error updating role:", err.Error())
		return nil, errors.New("error updating role")
	}
	return &role, nil
}

/* 
GetRole retrieves a single role by ID
@params roleID
*/

func GetRole(roleID uuid.UUID) (*Role, error) {
	
	var role Role
	if err := db.First(&role, "id = ?", roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("role not found:", err.Error())
			return nil, errors.New("role not found")
		}
		log.Println("error retrieving role:", err.Error())
		return nil, errors.New("error retrieving role")
	}
	return &role, nil
}

/*
GetAllRoles retrieves all roles from the database

*/
func GetAllRoles() ([]Role, error) {
	var roles []Role
	if err := db.Find(&roles).Error; err != nil {
		log.Println("error retrieving roles:", err.Error())
		return nil, errors.New("error retrieving roles")
	}
	return roles, nil
}

// DeleteRole deletes a role from the database by ID
func DeleteRole(roleID uuid.UUID) error {
	// Check if the role exists first
	var role Role
	if err := db.First(&role, "id = ?", roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("role not found")
			return errors.New("role not found")
		}
		log.Println("error finding role:", err.Error())
		return errors.New("error finding role")
	}

	// Role exists, proceed to delete
	if err := db.Delete(&role).Error; err != nil {
		log.Println("error deleting role:", err.Error())
		return errors.New("error deleting role")
	}

	return nil
}