package main

import "net/http"

// DummyAuthorizer is a concrete implementation of the Authorizer interface,
// where all http requests are authorized.
type DummyAuthorizer struct{}

// NewDummyAuthorizer creates a new dummy authorizer.
func NewDummyAuthorizer() *DummyAuthorizer {
	return &DummyAuthorizer{}
}

func (d *DummyAuthorizer) Authorize(r *http.Request) (bool, error) {
	// Always authorize call
	return true, nil
}
