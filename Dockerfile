# Start from a minimal base image with Go installed
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port on which the application will listen
EXPOSE 8080

# Run the Go application
CMD ["./main"]
