package services

import (
	"errors"

	"github.com/dancankarani/medicare/internal/app/database"
	"github.com/dancankarani/medicare/internal/app/model"
	"gorm.io/gorm"
)

var db = database.ConnectDB()

func CheckUserRegistered(email, phoneNumber string) (map[string][]string, error) {
    var user model.User
    errorsMap := make(map[string][]string)

    // Check if email exists
    err := db.Where("email = ?", email).First(&user).Error
    if err == nil {
        errorsMap["email"] = append(errorsMap["email"], "User with the given email is already registered.")
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errors.New("database error while checking email")
    }

    // Check if phone number exists
    err = db.Where("phone_number = ?", phoneNumber).First(&user).Error
    if err == nil {
        errorsMap["phone"] = append(errorsMap["phone"], "User with the given phone number is already registered.")
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errors.New("database error while checking phone number")
    }

    // If errors exist in the map, return them
    if len(errorsMap) > 0 {
        return errorsMap, errors.New("validation failed")
    }

    // No errors, return nil
    return nil, nil
}
