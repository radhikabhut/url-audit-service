# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o url-audit-service ./cmd/server/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/url-audit-service .
# Copy config
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./url-audit-service"]
