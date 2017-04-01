package main

import (
	"fmt"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	s := newServer()
	fmt.Print("Server started and listening on ", port)
	s.Run(":" + port)
}
