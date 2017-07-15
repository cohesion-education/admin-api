package main

import (
	"os"

	"github.com/cohesion-education/api/pkg/cohesioned/http"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3001"
	}

	http.Run(port)
}
