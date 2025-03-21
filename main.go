package main

import (
	"fmt"
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/api/routes"
)

func main() {
	model.DbMigrator()
	fmt.Println("hello medicare")
	routes.RegisterEndpoints()
}
