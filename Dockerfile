# Build stage
FROM golang:1.22.2-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o k8s-resource-analyzer-api ./cmd/api

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/k8s-resource-analyzer-api .

# Copy .env.example as default config
COPY .env.example .env

# Expose port
EXPOSE 9000

# Run the application
CMD ["./k8s-resource-analyzer-api"] 