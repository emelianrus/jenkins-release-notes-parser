# First stage: build the Go app
FROM golang:1.17-alpine AS build

# Set the working directory to /app
WORKDIR /app

# Copy all files except frontend directory into the container at /app
COPY . /app

# Download dependencies
RUN go mod download

# Build the Go app
RUN go build -o main .

# Second stage: create a lightweight container to run the app
FROM alpine:latest

# Set the working directory to /app
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/main .

# Expose port 8080 for the app to listen on
EXPOSE 8080

# Set the entry point to run the binary
ENTRYPOINT ["./main"]