package main

import (
	"io"
	"time"
)

// Storage wraps the interface for backend storage engine.
//
// PutObject prepares an object that can be written on the underlying engine.
// The object is unikely identified by the key.
//
// GetObject gets an object from storage identified by the key. If the object
// doesn't exist on the underlying engine, it should return nil.
type Storage interface {
	PutObject(key string) StorageObject
	GetObject(key string) StorageObject
}

// StorageObject wraps the interface for a single object on the storage egine.
//
// The object should be writable, readable and seekable.
type StorageObject interface {
	io.Writer
	io.Closer
	io.Reader
	io.Seeker

	// Clean should be called when the transfer is interruped, to it gives the
	// oportunity to the backend storage to clean after itself
	Clean() error

	// Name should return the name (key) of the file
	Name() string

	// Modtime should return the last modification of the file
	Modtime() time.Time
}
