# Start with a base Go image
FROM golang:1.20.5-alpine

# Set environment variables
ENV CGO_ENABLED=1
ENV GOOS=linux

# Install any necessary system dependencies
RUN apk add --update --no-cache gcc libc-dev

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

EXPOSE 8080

# Build the Go application
RUN go build -o main

# Run the database migrations
# RUN go run migration/migration.go

# Set the entrypoint command
CMD ["./main"]
