package services

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"gorm.io/gorm"
)

// User struct for illustration
type User struct {
    Username string
    Role     string
}
// GenerateRoleBasedUsername generates a username based on the role and ensures it’s unique
func GenerateRoleBasedUsername(role string) (string, error) {
    maxAttempts := 10
    var username string
    rand.Seed(time.Now().UnixNano()) // Seed the random number generator

    for i := 0; i < maxAttempts; i++ {
        // Generate a candidate username based on role
        randomCode := fmt.Sprintf("AB%d", rand.Intn(9000)+1000) // Random code in range AB1000-AB9999
        
        switch strings.ToLower(role) {
        case "admin":
            username = fmt.Sprintf("Admin@%s", randomCode)
        case "doctor":
            username = fmt.Sprintf("Doctor@%s", randomCode)
        // Add other roles as needed
        default:
            username = fmt.Sprintf("Staff@%s", randomCode)
        }

        // Check if the username already exists
        var user User
        err := db.Where("username = ?", username).First(&user).Error

        if err == gorm.ErrRecordNotFound {
            // Username is unique
            return username, nil
        } else if err != nil {
            // Handle other errors from the database
            return "", fmt.Errorf("database error: %v", err)
        }
    }

    // Fallback if a unique username couldn’t be generated within the max attempts
    fallbackUsername := fmt.Sprintf("%s@%s%d", strings.Title(role), "AB", time.Now().Unix()%10000)
    return fallbackUsername, nil
}
