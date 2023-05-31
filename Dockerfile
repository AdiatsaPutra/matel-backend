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
RUN go mod download

# Copy the rest of the application source code
COPY . .

# # Build the Go application
# RUN go build -o matel-backend
# RUN go run main.go

ADD . .

EXPOSE 8000

ENTRYPOINT CompileDaemon --build-"go build main.go"
--command=./main

RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon


# # Set the entrypoint command
# CMD ["./matel-backend"]