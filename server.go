package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Server struct {
	port        int
	timeout     time.Duration
	maxBodySize int64
	authorizer  Authorizer
	storage     Storage
}

// Start starts the http server and blocks forever.
func (s *Server) Start() {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.port),
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           s,
	}

	log.Fatal(srv.ListenAndServe())
}

// authorizeHandler is a middleware to check for proper authorization
func (s *Server) authorizeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok, err := s.authorizer.Authorize(r)

		if err != nil {
			http.Error(w, "Internal authorization error", http.StatusInternalServerError)
			return
		}

		if !ok {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// checkKeyPresence is a middleware to check for required key path
func (s *Server) checkKeyPresence(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check path presence
		if r.URL.Path == "" {
			http.Error(w, "Missing object key", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleUpload is responsibile for handling an upload request
func (s *Server) handleUpload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that we have something to upload
		if r.Body == nil {
			http.Error(w, "Empty body", http.StatusBadRequest)
			return
		}

		// Check request size, if possible
		if r.ContentLength > 0 && r.ContentLength > s.maxBodySize {
			http.Error(w, "Request too big", http.StatusRequestEntityTooLarge)
			return
		}

		// Clients might still cheat, so create a reader that doesn't read more than allowed.
		body := http.MaxBytesReader(w, r.Body, s.maxBodySize)

		// Ask the storage backend for an object
		object := s.storage.PutObject(r.URL.Path)
		if object == nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Pump bytes from http to storage
		_, err := io.Copy(object, body)
		if err != nil {
			object.Clean()

			if err == http.ErrHandlerTimeout {
				// net.TimeoutHandler already returns the proper error message to the client
			} else {
				log.Printf("Error pumping bytes: %s", err.Error())
				http.Error(w, "Storage error", http.StatusInternalServerError)
			}

			return
		}

		// Tell the storage backend that we are done with the transfer.
		err = object.Close()
		if err != nil {
			log.Printf("Error finalizing file: %s", err.Error())
			http.Error(w, "Storage error", http.StatusInternalServerError)
			return
		}

		// If we get here, the upload succeeded, so just tell the client everything's OK
		w.WriteHeader(200)
		_, _ = w.Write([]byte("OK"))
	})
}

// handleDownload is responsibile for handling a download request
func (s *Server) handleDownload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ask storage for file
		object := s.storage.GetObject(r.URL.Path)
		if object == nil {
			http.NotFound(w, r)
			return
		}

		// Stream file to client directly
		http.ServeContent(w, r, object.Name(), object.Modtime(), object)
	})
}

// ServeHTTP is the main handler where all the requests go.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler = func(next func() http.Handler) {
		http.TimeoutHandler( // Ensures the request has a timeout
			s.authorizeHandler( // Ensures the request is authorized
				s.checkKeyPresence( // Ensures the request has a required key path
					next(),
				),
			),
			s.timeout,
			http.StatusText(http.StatusGatewayTimeout),
		).ServeHTTP(w, r)
	}

	switch r.Method {
	case "POST":
		handler(s.handleUpload)
	case "GET":
		handler(s.handleDownload)
	default:
		http.Error(w, "Unknown method", http.StatusBadRequest)
	}
}
