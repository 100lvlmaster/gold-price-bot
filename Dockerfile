# Stage 1: Build the binary
FROM golang:1.25-alpine AS builder

# Install build dependencies for CGO (required by sqlite3)
RUN apk add --no-cache build-base

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go application
# CGO_ENABLED=1 is required for github.com/mattn/go-sqlite3
# -w -s flags strip debug symbols and DWARF tables to minimize binary size
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o gold-price-api ./cmd/api/main.go

# Stage 2: Final minimal production image
FROM alpine:latest

# Install CA certificates (required for secure Telegram API communication)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the pre-built binary
COPY --from=builder /app/gold-price-api .

# Copy .env
COPY --from=builder /app/.env .

# Expose the API port
EXPOSE 8080

# Run the binary
CMD ["./gold-price-api"]
