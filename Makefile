BINARY_NAME_CLIENT=client
BINARY_NAME_SERVER=server

dev:
	@echo "-->Run dev mode"
	go run ./cmd/$(BINARY_NAME).go

test:
	echo "-->Run test"
	go test -v ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

clean:
	@echo "-->Clean"
	go clean
	rm -rf test.db
	rm -rf bin

dep: clean
	@echo "-->Download dependencies"
	go mod download
	go mod verify
	go mod tidy

build_sever: dep
	@echo "==>Building binary"
	go build -o bin/ -v ./cmd/server/$(BINARY_NAME_SERVER).go

run_server: build_server
	@echo "==>Run binary"
	./bin/$(BINARY_NAME_SERVER)

build_client: dep
	@echo "==>Building binary for client"
	go build -o bin/ -v ./cmd/client/$(BINARY_NAME_CLIENT).go

run_client: build_client
	@echo "==>Run binary for client"
	./bin/$(BINARY_NAME_CLIENT)