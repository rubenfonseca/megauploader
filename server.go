package main

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	port int
}

func (s *Server) Start() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s,
	}

	log.Fatal(srv.ListenAndServe())
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, _ = w.Write([]byte("Hello world"))
}
