package main

import (
	"log"
)

func main() {
	server := &Server{
		port: 9292,
	}

	log.Printf("Starting server on port %d", server.port)
	server.Start()
}
