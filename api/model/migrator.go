package model

import "fmt"

func DbMigrator() {
	fmt.Println("initializing db migrator")
	db.AutoMigrate(
		&ClinicBooking{},
		&User{},
		&Billing{},
		&Payment{},
		&LabTest{},
		&Patient{},
		&Medicine{},
		&Inventory{},
		&Referral{},
		&Prescription{},
	)
}