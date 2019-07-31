package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"
)

func TestMyHandler(t *testing.T) {
	root, _ := ioutil.TempDir("", "megaupload")

	storage := NewFileStorage()
	storage.Root = root

	handler := &Server{
		storage:    storage,
		authorizer: NewDummyAuthorizer(),
		timeout:    10 * time.Second,
	}

	t.Run("GET key that doesn't exist", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/foo", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		if response.Code != 404 {
			t.Error("Expected 404 response, got", response.Code)
		}
	})

	t.Run("GET key that exists", func(t *testing.T) {
		// Create dummy file
		tempfile, _ := ioutil.TempFile(root, "dummy")
		ioutil.WriteFile(tempfile.Name(), []byte("1 2 3"), 0644)

		path, _ := filepath.Rel(root, tempfile.Name())

		request, _ := http.NewRequest(http.MethodGet, "/"+path, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		if response.Code != 200 {
			t.Error("Expected 200 response, got", response.Code)
		}

		if response.Body.String() != "1 2 3" {
			t.Error("Expecting body 1 2 3, got", response.Body.String())
		}
	})
}
