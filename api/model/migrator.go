package model

import "fmt"

func DbMigrator() {
	fmt.Println("initializing db migrator")
	db.AutoMigrate(
		&User{},
		&Patient{},
		&Inventory{},
		&Medicine{},
		Prescription{},
	)
}