FROM golang:1.23-alpine AS builder
ENV GO111MODULE=on
RUN mkdir /tmpdir
WORKDIR  /tmpdir

# Install golangci-lint
RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.5

COPY . .
RUN go mod download

# Run linter
RUN golangci-lint run --timeout=5m

RUN CGO_ENABLED=0 go build -o company

# Defining App image
FROM alpine:latest
RUN apk add --no-cache --update ca-certificates

WORKDIR /app
# Copy App binary to image
COPY --from=builder /tmpdir/company /app/
EXPOSE 8080

ENTRYPOINT ["./company"]
