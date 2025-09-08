# Use a multi-stage build to keep the final image small
# Stage 1: The builder stage
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install Air for hot-reloading
RUN go install github.com/air-verse/air@latest

# Copy Go module files to cache dependencies
COPY go.mod go.sum ./

# Download all project dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the binary to ensure dependencies are resolved
RUN go build -o /tmp/auth_service ./cmd/main.go

# Stage 2: The final runtime image
FROM golang:1.25-alpine

# Copy the Air binary and the application source code from the builder stage
COPY --from=builder /go/bin/air /usr/local/bin/air
COPY --from=builder /app /app

# Set the working directory
WORKDIR /app

# Set Air as the entrypoint for hot-reloading
ENTRYPOINT ["air"]