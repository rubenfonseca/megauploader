package main

import (
	"io"
)

// Storage wraps the interface for backend storage engine.
//
// PutObject prepares an object that can be written on the underlying engine.
// The object is unikely identified by the key.
type Storage interface {
	PutObject(key string) StorageObject
}

// StorageObject wraps the interface for a single object on the storage egine.
//
// The object should be writable, readable and seekable.
type StorageObject interface {
	io.Writer
	io.Closer
}
