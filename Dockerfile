# Start with a base Go image
FROM golang:1.20

# Set environment variables
ENV CGO_ENABLED=1
ENV GOOS=linux

# Install any necessary system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go get github.com/githubnemo/CompileDaemon
RUN go mod download

# Copy the rest of the application source code
COPY . .

EXPOSE 8000
# # Build the Go application


RUN go build -o main
RUN go run main.go

# # Run the Go application
CMD CompileDaemon -log-prefix=false -build="go build -o main" -command="./main"


# CMD ["./matel-backend"]