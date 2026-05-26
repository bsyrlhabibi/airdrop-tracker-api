.PHONY: run build dev clean swag test

run:
	@echo "Starting server..."
	go run cmd/server/main.go

build:
	@echo "Building binary..."
	go build -o bin/server cmd/server/main.go

dev:
	@echo "Starting server..."
	go run cmd/server/main.go

clean:
	rm -rf bin/ data/*.db

swag:
	@echo "Generating Swagger docs..."
	swag init -g cmd/server/main.go

test:
	go test ./... -v
