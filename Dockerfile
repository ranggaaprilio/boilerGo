# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o boilergo .

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata mysql-client make wget go git

# Copy the binary from builder
COPY --from=builder /app/boilergo .
COPY --from=builder /app/config.docker.yml /app/config.yml

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./boilergo"]
