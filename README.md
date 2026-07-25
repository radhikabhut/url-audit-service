# URL Audit Service

A production-grade, highly resilient microservice written in Go designed to inspect target URLs, measure performance, parse HTML titles safely (SSRF protected), and return structured metadata.

---

## Architecture

The project follows the principles of **Clean Architecture** to ensure separations of concern, high testability, and decoupling of business rules from concrete framework integrations.

```
cmd/
  server/             - Application entrypoint & dependency injection bootstrapping
internal/
  config/             - Viper configurations loaded via .env file & env overrides
  logger/             - Structured JSON log configurations using slog
  middleware/         - RequestID, recovery, CORS and rate-limiting middlewares
  validator/          - Go-Playground validator integrations & SSRF URL safety checker
  cache/              - Thread-safe local memory cache with configurable TTL and normalizer
  limiter/            - Thread-safe per-IP rate limiter
  client/             - Outbound HTTP audit client with timeout control & html parser
  repository/         - In-memory thread-safe repository storing audit histories
  service/            - Service layer orchestrating domain logic (caching + semaphore concurrency limits)
  handler/            - Gin HTTP controller layer delivering normalized error payloads
  response/           - API response schemas (success data wrapper + structured errors)
  errors/             - Domain error classifications
```

---

## Setup & Running Locally

### Prerequisites
- Go 1.25+ (Standard installation)

### 1. Copy Environment Configuration
Verify that the `.env` configuration file exists at the root of the project:

```env
PORT=8080
REQUEST_TIMEOUT=10s
CACHE_TTL=5m
MAX_CONCURRENT_AUDITS=100
RATE_LIMIT=10
RATE_LIMIT_WINDOW=1s
HOST=0.0.0.0
LOG_LEVEL=info
LOG_FORMAT=json
```

### 2. Launch the Application
Run the main server directly using Go:

```bash
go run cmd/server/main.go
```

---

## Environment Variables

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PORT` | Listening socket port for the HTTP API | `8080` |
| `HOST` | Listening interface binding host | `0.0.0.0` |
| `REQUEST_TIMEOUT` | Timeout limit for outbound url audits | `10s` |
| `CACHE_TTL` | Retention expiration time for cached results | `5m` |
| `MAX_CONCURRENT_AUDITS` | Max concurrent outbound HTTP connection count | `100` |
| `RATE_LIMIT` | Token bucket size for per-client IP throttling | `10` |
| `RATE_LIMIT_WINDOW` | Time window for rate limiter token refill | `1s` |
| `LOG_LEVEL` | Logging level (`debug`, `info`, `warn`, `error`) | `info` |
| `LOG_FORMAT` | Logger output format (`json`, `text`) | `json` |

---

## API Documentation

### 1. Health Check
Checks the status of the server.

- **URL**: `/health`
- **Method**: `GET`
- **Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "status": "UP"
  }
}
```

### 2. Audit URL
Enforces SSRF prevention checks (blocks loopbacks, private IPs, and localhost) and audits the target URL.

- **URL**: `/api/v1/audit`
- **Method**: `POST`
- **Headers**: `Content-Type: application/json`
- **Request Body**:
```json
{
  "url": "https://example.com"
}
```
- **Response (200 OK)**:
```json
{
  "url": "https://example.com",
  "reachable": true,
  "statusCode": 200,
  "responseTimeMs": 143,
  "contentType": "text/html",
  "contentLength": 1250,
  "title": "Example Domain",
  "cached": false,
  "checkedAt": "2026-07-25T14:40:32Z"
}
```

- **Response (400 Bad Request - SSRF blocked)**:
```json
{
  "requestId": "54bfa136-fb57-4184-be71-b60a5094c7fe",
  "error": {
    "code": "INVALID_URL",
    "message": "connections to loopback addresses are rejected"
  }
}
```

- **Response (429 Too Many Requests - Throttled)**:
```json
{
  "requestId": "4c9edf25-1246-419f-adea-11eb9efef399",
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "rate limit exceeded"
  }
}
```

---

## Testing

Execute formatting checks, static analysis code vetting, and the complete unit/integration test suite:

```bash
# Verify formatting rules
gofmt -l .

# Run static analyzer
go vet ./...

# Run complete test suite with coverage
go test -v -coverprofile=coverage.out ./...

# View statement coverage details
go tool cover -func=coverage.out
```

---

## Deployment

Build and run the application containerized via Docker and Docker Compose:

```bash
# Build and run using Docker Compose
docker compose up --build -d

# Stop and clean up containers
docker compose down
```
