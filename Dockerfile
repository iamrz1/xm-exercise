FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /company-service ./cmd/server

# Use a small alpine image
FROM alpine:3.18

WORKDIR /

# Copy the binary from builder
COPY --from=builder /company-service /company-service

# Create a non-root user
RUN adduser -D -g '' appuser
USER appuser

# Run the application
ENTRYPOINT ["/company-service"] 