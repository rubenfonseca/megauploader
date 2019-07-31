# MegaUploader

Simple example of a file upload service built in Go. The objective is to build the
server just using Go's stdlib.

## How to use

### Using go

`make run`

### Using docker

`make docker`

## Implemented features

- Sane defaults
- Streaming architecture (don't exhaust memory)
- Pluggable backend storage engine
- Pluggable authorization engine

## Future features

- Pluggable metadata storage engine
