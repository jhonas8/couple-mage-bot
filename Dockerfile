# Use the official Golang image to create a build artifact.
# This image is based on Debian, so we use apt-get to install packages.
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o main .

# Use a Docker multi-stage build to create a lean production image.
# Start with a smaller base image
FROM debian:buster-slim

# Set the Current Working Directory inside the container
WORKDIR /app

RUN apt-get update -y && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    build-essential

RUN apt-get update -y
RUN apt-get install -y libx11-dev

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]

