# Build stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /workspace

# Copy source code and vendored dependencies
COPY . .

# Build the application using vendored modules
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ayame cmd/ayame/main.go

# Final stage
FROM alpine:3.20

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /workspace/ayame .

# Copy the example config file
COPY config_example.ini ./config.ini

# Create logs directory
RUN mkdir -p logs

# Expose the default port
EXPOSE 3000

# Expose the default Prometheus port
EXPOSE 4000

# Run the binary
CMD ["./ayame", "-C", "/app/config.ini"]