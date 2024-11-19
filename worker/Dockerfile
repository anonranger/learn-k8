# Use an official Golang image as the base image
FROM golang:1.23.0-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o backend .

# Use a minimal base image for the final stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/backend .

# Expose the port the app runs on
EXPOSE 80

# Set environment variables for MongoDB URI and Port
ENV MONGO_URI=mongodb://localhost:27017/
ENV PORT=3000

# Command to run the application
CMD ["./backend"]