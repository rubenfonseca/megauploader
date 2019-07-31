.PHONY: clean build run docker
.DEFAULT_GOAL := build

clean:
	rm -f ./megauplaoder

build:
	go build -o megauploader .

run: 
	go run .

docker:
	docker build . -t megauploader
	docker run --rm -it -p 9292:9292 megauploader
