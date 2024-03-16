# Use the official Golang image to create a build artifact.
FROM golang:1.20 AS builder

# Create a working directory inside the container.
WORKDIR /app

# Copy the Go modules files.
COPY go.mod ./

# Download the Go dependencies.
RUN go mod download

# Copy the source code into the container.
COPY . .

# Set the CGO_ENABLED environment variable to 0.
ENV CGO_ENABLED=0

# Build the Go application.
RUN go build -o blitz .

# Use a minimal base image to reduce the final container size.
FROM alpine:latest as runner

# Copy the binary from the builder stage to the new image.
COPY --from=builder /app/blitz /usr/local/bin/blitz

# Set the working directory.
WORKDIR /app

# Command to run the Go application.
ENTRYPOINT [ "blitz"]
