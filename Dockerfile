# Use the official Golang image to create a build artifact.
# This image is based on Debian, so we use apt-get to install packages.
FROM golang:1.22 as builder

# Set the working directory inside the container
WORKDIR /goapp

# Copy the Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

EXPOSE 8080

# Set the entry point command to run the built binary
CMD ["./main"]
