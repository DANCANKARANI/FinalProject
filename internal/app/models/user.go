package models

import "time"

type User struct {
	ID          string `json:"id" gorm:"type:varchar(36)"`
	FullName    string `json:"full_name" gorm:"type:varchar(100)"`
	Email       string `json:"email" gorm:"type:varchar(100)"`
	PhoneNumber string `json:"phone_number" gorm:"type:varchar(15)"`
	DateOfBirth time.Time	`json:"date_of_birth"`
	Password	string	`json:"password" gorm:"type:varchar(255)"`
	Adrress		string	`json:"address" gorm:"type:text"`
	Gender		string	`json:"gender" gorm:"type:varchar(10)"`
	CreatedAt	time.Time	`json:"created_at" `
}