package main

import (
	"log"
)

func main() {
	server := &Server{
		port:        9292,
		maxBodySize: 1 * 1024 * 1024 * 1024, // 1GB
		authorizer:  NewDummyAuthorizer(),
		storage:     NewFileStorage(),
	}

	log.Printf("Starting server on port %d", server.port)
	server.Start()
}
