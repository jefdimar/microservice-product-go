# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install required build tools
RUN apk add --no-cache git

# Copy only the files needed for go mod download
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .

# Final stage - using scratch (smallest possible image)
FROM scratch

# Copy SSL certificates for HTTPS requests
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

# Copy only the necessary files from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]