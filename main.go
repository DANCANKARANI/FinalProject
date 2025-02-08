package main

import (
	"fmt"

	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/api/routes"
)

func main() {
	fmt.Println("hello medicare")
	model.DbMigrator()
	routes.RegisterEndpoints()
	
}