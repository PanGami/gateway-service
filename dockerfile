# Use a Golang base image
FROM golang:1.22.0

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Install Go dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o main .

# Expose the API port
EXPOSE 7080

# Command to run the application
CMD ["./main"]
