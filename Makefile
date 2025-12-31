# Makefile for building and running Go server and client applications

build:
	go build -o bin/main ./main.go
	sudo cp bin/main /usr/local/bin/bloc

# Build server binary
build-server:
	go build -o bin/server ./server

# Build client binary
build-client:
	go build -o bin/client ./client

build-bloc:
	go build -o bin/bloc ./cmd

# Run server (auto-build)
run-server: build-server
	./bin/server

# Run client (auto-build)
run-client: build-client
	./bin/client

run-bloc: build-bloc
	./bin/bloc

run:
	make build
	./bin/main

# Clean binaries
clean:
	rm -rf bin
