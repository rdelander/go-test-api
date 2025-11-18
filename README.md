# Go Test API

A production-ready REST API built with Go following enterprise best practices, featuring clean architecture, dependency injection, and comprehensive testing.

## Project Structure

```
go-test-api/
├── cmd/
│   └── server/          # Application entry point
│       └── main.go
├── internal/            # Private application code
│   ├── handler/         # HTTP request handlers
│   │   ├── hello.go
│   │   └── hello_test.go
│   ├── model/           # Data models and DTOs
│   │   └── hello.go
│   ├── server/          # Server setup and routing
│   │   └── server.go
│   └── validator/       # Validation logic
│       └── validator.go
├── pkg/                 # Public, reusable packages
│   └── response/        # HTTP response helpers
│       └── response.go
├── go.mod
└── go.sum
```

## Requirements

- Go 1.22 or higher

## Setup

1. Clone the repository
2. Navigate to the project directory:
   ```bash
   cd go-test-api
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```

## Running the Server

### Using go run:

```bash
go run cmd/server/main.go
```

### Using air (auto-reload for development):

```bash
air
```

The server will start on `http://localhost:8080`

## API Endpoints

### GET /hello_world

Returns a default greeting message.

**Request:**

```bash
curl http://localhost:8080/hello_world
```

**Response:**

```json
{
  "message": "Hello, World!"
}
```

### POST /hello_world

Accepts a JSON payload with a name and returns a personalized greeting message.

**Request:**

```bash
curl -X POST http://localhost:8080/hello_world \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice"}'
```

**Request Body:**

```json
{
  "name": "Alice"
}
```

**Response:**

```json
{
  "message": "Hello, Alice!"
}
```

**Validation:**

- `name` is required
- `name` must be between 1-100 characters

**Error Response:**

```json
{
  "error": "Field 'Name' failed validation 'required'"
}
```

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test ./... -cover
```

Run tests with verbose output:

```bash
go test ./... -v
```

## Building

To build a binary:

```bash
go build -o api-server cmd/server/main.go
./api-server
```

## Architecture

This project follows **Clean Architecture** principles:

- **cmd/**: Application entry points (main packages)
- **internal/**: Private application code that cannot be imported by other projects
  - **handler/**: HTTP handlers - thin layer that handles HTTP concerns
  - **model/**: Domain models and data transfer objects
  - **server/**: Server initialization and routing configuration
  - **validator/**: Business validation logic
- **pkg/**: Public libraries that can be imported by external projects

### Dependency Injection

All dependencies are injected through constructors, making the code:

- Easier to test (can inject mocks)
- More maintainable (clear dependencies)
- More flexible (easy to swap implementations)

### Adding New Endpoints

1. Define models in `internal/model/`
2. Create handler in `internal/handler/`
3. Write tests in `internal/handler/*_test.go`
4. Register route in `internal/server/server.go`
