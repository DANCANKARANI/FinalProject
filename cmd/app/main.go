package main

import (
	"fmt"
	"github.com/dancankarani/medicare/internal/endpoints"
)

func main() {
	fmt.Println("hello medicare")
	endpoints.RegisterEndPoints()
}