# ==========================================
# Stage 1: Build Stage
# ==========================================
FROM golang:1.25-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy dependency file first (To optimize the Docker cache)
COPY go.mod go.sum ./

# Download all dependency
RUN go mod download

# Copy all source code to the container
COPY . .

# Build the application into a binary named 'main'
# Ensure that the path ‘cmd/api/main.go’ matches your folder structure
RUN go build -o main cmd/api/main.go

# ==========================================
# Stage 2: Run Stage
# ==========================================
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install CA Certificates (It is important if the application requests HTTPS externally)
RUN apk --no-cache add ca-certificates

# Copy binary 'main' from Stage 1
COPY --from=builder /app/main .

# Copy configuration file (app.env) and migration folder (if it needs to be run inside)
COPY .env .
COPY db/migrations ./db/migrations

# Expose app port (according to your configuration, for example 3000)
EXPOSE 3000

# Command to run the application
CMD ["./main"]