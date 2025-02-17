# ========== STAGE 1: Build the Go binary ==========
FROM golang:1.20-alpine AS build

# Set your working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (for dependency caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
# Example: if your main is in cmd/main.go, you might do:
RUN go build -o myapp ./cmd

# ========== STAGE 2: Create a minimal runtime image ==========
FROM alpine:3.17

# Create a non-root user (optional, but recommended for security)
RUN addgroup -S mygroup && adduser -S myuser -G mygroup

# Copy the binary from the builder stage
COPY --from=build /app/myapp /usr/local/bin/myapp

# Make sure the binary has execution permissions (usually by default, but to be sure):
RUN chmod +x /usr/local/bin/myapp

# Switch to the non-root user
USER myuser

# Expose a port if your app listens on it, e.g., 8080
EXPOSE 8080

# Command to run your binary
ENTRYPOINT ["myapp"]