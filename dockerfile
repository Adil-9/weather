# Use an official Go runtime as a parent image
FROM golang:alpine

# Set the working directory in the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application
RUN go get -d -v ./... && \
    go install -v ./...

CMD ["go", "run", "cmd/main.go"]
