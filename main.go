package main

import (
	"log"
)

func main() {
	server := &Server{
		port:       9292,
		authorizer: NewDummyAuthorizer(),
	}

	log.Printf("Starting server on port %d", server.port)
	server.Start()
}
