# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy gomod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Run Stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies (check if we need ca-certificates for ssl)
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/config/config.yaml ./config/config.yaml
# COPY --from=builder /app/migrations ./migrations # Uncomment if you have migrations

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]
