package model

import "fmt"

func DbMigrator() {
	fmt.Println("initializing db migrator")
	db.AutoMigrate(
		&User{},
		&LabTest{},
		&Patient{},
		&Medicine{},
		&Inventory{},
		&Referral{},
		&Prescription{},
	)
}