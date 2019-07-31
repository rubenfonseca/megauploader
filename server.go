package main

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	port        int
	maxBodySize int64
	authorizer  Authorizer
}

func (s *Server) Start() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s,
	}

	log.Fatal(srv.ListenAndServe())
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Authorization
	ok, err := s.authorizer.Authorize(r)
	if err != nil {
		http.Error(w, "Internal authorization error", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	// Check path presence
	if r.URL.Path == "" {
		http.Error(w, "Missing object key", http.StatusBadRequest)
		return
	}

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
	_ = http.MaxBytesReader(w, r.Body, s.maxBodySize)

	w.WriteHeader(200)
	_, _ = w.Write([]byte("Hello world"))
}
