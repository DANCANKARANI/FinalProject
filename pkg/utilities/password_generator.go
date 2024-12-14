package utilities

import (
	"fmt"
	"math/rand"
	"time"
)

func PasswordGenerator() string {
	// Get the current year
	currentYear := time.Now().Year()

	// Generate 4 random numbers
	rand.Seed(time.Now().UnixNano())
	randomNumbers := rand.Intn(10000) // Generates a number between 0 and 9999

	// Format the password
	password := fmt.Sprintf("%d.%04d", currentYear, randomNumbers)

	return password
}