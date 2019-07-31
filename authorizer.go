package main

import "net/http"

// Authorizer is the interface what wraps the basic authorization engine.
//
// Authorize receives a http request and decides if the request is
// authenticated and authorized to proceed.
//
// You can create your own authorizer concrete class that may connect to
// database, check JWT tokens, cookies, etc, and decide if the request is
// authorized or not. Having access to the original http request gives you
// access to all the headers.
type Authorizer interface {
	Authorize(r *http.Request) (bool, error)
}
