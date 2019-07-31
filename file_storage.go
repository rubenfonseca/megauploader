package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

// FileStorage is a concrete implementation of the Storage interface, where the
// local filesystem is used as a backing storage.
//
// The prefix where the files are stored is saved on Root.
type FileStorage struct {
	Root string
}

// FileStorageObject is a concrete implementatin of the Storage interface. The
// object is backed by a local os.File.
type FileStorageObject struct {
	file *os.File
	path string
}

// NewFileStorage creates a new FileStorage with prefix on /tmp
func NewFileStorage() *FileStorage {
	return &FileStorage{
		Root: "/tmp",
	}
}

func (f *FileStorage) PutObject(key string) StorageObject {
	path := f.buildPath(key)

	file, err := os.Create(path)
	if err != nil {
		log.Printf("Error creating object on file storage: %s", err.Error())
		return nil
	}

	return &FileStorageObject{
		file: file,
		path: path,
	}
}

func (f *FileStorage) GetObject(key string) StorageObject {
	path := f.buildPath(key)

	// Check if file exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}

	if info.IsDir() {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil
	}

	return &FileStorageObject{
		file: file,
		path: path,
	}
}

func (f *FileStorage) buildPath(key string) string {
	return filepath.Join(f.Root, filepath.FromSlash(path.Clean("/"+key)))
}

func (f *FileStorageObject) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

func (f *FileStorageObject) Close() error {
	return f.file.Close()
}

func (f *FileStorageObject) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func (f *FileStorageObject) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}
