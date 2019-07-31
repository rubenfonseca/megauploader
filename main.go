package main

import (
	"log"
	"time"
)

func main() {
	// Initialize the server with some defaults
	server := &Server{
		port:        9292,
		timeout:     5 * time.Minute,
		maxBodySize: 1 * 1024 * 1024 * 1024, // 1GB
		authorizer:  NewDummyAuthorizer(),
		storage:     NewFileStorage(),
	}

	log.Printf("Starting server on port %d", server.port)
	server.Start()
}
