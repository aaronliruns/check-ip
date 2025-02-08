FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o check-ip

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/check-ip .

# Copy configuration files
COPY config.yaml .
COPY bsgsth.txt .

# Create non-root user
RUN adduser -D appuser
USER appuser

# Set production mode for Gin
ENV GIN_MODE=release

EXPOSE 8888

CMD ["./check-ip"]
