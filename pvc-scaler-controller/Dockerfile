# Use the official Go image as the base image
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o pvc-scaler ./cmd/main.go

# Use a minimal image for the runtime
FROM alpine:3.18

# Install necessary packages
RUN apk add --no-cache ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/pvc-scaler .

# Command to run the application
CMD ["./pvc-scaler"]

