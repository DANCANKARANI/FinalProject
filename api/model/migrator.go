package model

import "fmt"

func DbMigrator() {
	fmt.Println("initializing db migrator")
	db.AutoMigrate(
		&Doctor{},
		&Patient{},
		&Pharmacist{},
		&Inventory{},
	)
}