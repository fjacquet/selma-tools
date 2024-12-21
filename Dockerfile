# Stage 1: Build the Go application
FROM golang:latest AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create and set the working directory
WORKDIR /app

# Copy Go module manifests and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Web UI binary
WORKDIR /app/web
RUN go build -o /app/web

# Stage 2: Create a minimal image for running the application
FROM alpine:latest
# Install dependencies
RUN apk --no-cache add ca-certificates && apk add --no-cach --update sudo

# Set working directory
WORKDIR /app
# add new user
RUN adduser -D default -h /app
USER default

# Copy the compiled binary from the builder stage
COPY --from=builder /app/web .

# Copy static files for the Web UI
COPY static ./static
COPY downloads ./downloads

# Expose the application port
EXPOSE 8080

# Run the Web UI application
CMD ["./web"]
