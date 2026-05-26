# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies (for SQLite CGO)
RUN apk add --no-cache gcc musl-dev

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -o server cmd/server/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/server .

# Copy swagger docs
COPY --from=builder /app/docs ./docs

# Create data directory for SQLite
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Run
CMD ["./server"]
