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
	timeOut     time.Duration
	maxBodySize int64
	authorizer  Authorizer
	storage     Storage
}

func (s *Server) Start() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s,
	}

	log.Fatal(srv.ListenAndServe())
}

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
			log.Printf("Error pumping bytes: %s", err.Error())
			http.Error(w, "Storage error", http.StatusInternalServerError)
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

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.TimeoutHandler(
		s.authorizeHandler(
			s.checkKeyPresence(
				s.handleUpload(),
			),
		),
		s.timeOut,
		http.StatusText(http.StatusGatewayTimeout),
	).ServeHTTP(w, r)
}
