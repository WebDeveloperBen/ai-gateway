# syntax=docker/dockerfile:1

# --- Build Stage ---
FROM golang:1.23-alpine AS builder

# Set build-time environment (only for this stage)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /src

# Copy go.mod and go.sum separately to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the application binary
RUN go build -trimpath -ldflags="-s -w" -o /main ./cmd/api/main.go

# --- Final Stage ---
FROM scratch

# Copy CA certificates for HTTPS support
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled Go binary
COPY --from=builder /main /main

# Use a non-root user for security (must manually create it in scratch)
# Define user and group IDs
USER 10001:10001

# Expose the HTTP port
EXPOSE 8080

# Set the binary as entrypoint
ENTRYPOINT ["/main"]
