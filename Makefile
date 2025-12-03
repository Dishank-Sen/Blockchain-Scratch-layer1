# Build server binary
build-server:
	go build -o bin/server ./server

# Build client binary
build-client:
	go build -o bin/client ./client

# Run server (auto-build)
run-server: build-server
	./bin/server

# Run client (auto-build)
run-client: build-client
	./bin/client

# Clean binaries
clean:
	rm -rf bin
