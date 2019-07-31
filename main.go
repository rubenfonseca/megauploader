package main

import (
	"log"
	"time"
)

func main() {
	server := &Server{
		port:        9292,
		timeOut:     5 * time.Minute,
		maxBodySize: 1 * 1024 * 1024 * 1024, // 1GB
		authorizer:  NewDummyAuthorizer(),
		storage:     NewFileStorage(),
	}

	log.Printf("Starting server on port %d", server.port)
	server.Start()
}
