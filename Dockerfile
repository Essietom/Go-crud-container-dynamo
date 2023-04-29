# Use the official Golang image as the base image
FROM golang:1.16.3-alpine3.13

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application binary inside the container
RUN go build -o main .

# Set the command to run the binary when the container starts
CMD ["./main"]

# Expose the port on which the application will listen
EXPOSE 8080