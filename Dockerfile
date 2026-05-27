FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/server .

RUN mkdir -p /app/data
VOLUME /app/data

ENV APP_PORT=8080
ENV DB_PATH=/app/data/airdrop.db
ENV GIN_MODE=release

EXPOSE 8080

CMD ["./server"]
