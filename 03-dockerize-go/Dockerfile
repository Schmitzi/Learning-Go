# Use latest debian stable version
FROM golang:1.25-alpine

# Define root dir
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Copy docs
COPY docs/ ./docs/

# Download Go dependencies
RUN go mod download

# Copy OpennSSL Keys
COPY key.pem cert.pem ./

# Copy the main file to /app
COPY main.go ./

# Build binary
RUN go build -o main .

# Expose the port the server uses
EXPOSE 8080

# Start server
CMD ["./main"]
