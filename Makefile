.PHONY: run build dev clean swag test deploy swag-install

run:
	@echo "🚀 Starting server..."
	go run cmd/server/main.go

build:
	@echo "📦 Building binary..."
	go build -o bin/server cmd/server/main.go

dev:
	@echo "🔄 Starting server (dev mode)..."
	go run cmd/server/main.go

clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/ data/*.db

swag:
	@echo "📝 Generating Swagger docs..."
	swag init -g cmd/server/main.go

swag-install:
	@echo "📥 Installing swag CLI..."
	go install github.com/swaggo/swag/cmd/swag@latest

test:
	go test ./... -v

deploy:
	@echo "🚀 Deploying to Fly.io..."
	flyctl deploy
