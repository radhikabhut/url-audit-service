# Stage 1: Build the React frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app

# Copy dependency files and install
COPY frontend/package*.json ./
RUN npm install

# Copy source and build
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.26-alpine AS backend-builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source code
COPY . .

# Copy compiled frontend static assets from Stage 1 into the backend's static directory
COPY --from=frontend-builder /app/dist ./static

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o url-audit-service ./cmd/server/main.go

# Stage 3: Runner
FROM alpine:latest

WORKDIR /app

# Copy Go binary, env settings, and static assets folder
COPY --from=backend-builder /app/url-audit-service .
COPY --from=backend-builder /app/.env .
COPY --from=backend-builder /app/static ./static

EXPOSE 8080

CMD ["./url-audit-service"]
