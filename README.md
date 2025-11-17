# Go Test API

A simple REST API built with Go that provides a POST endpoint for hello world functionality.

## Requirements

- Go 1.21 or higher

## Setup

1. Clone the repository
2. Navigate to the project directory:
   ```bash
   cd go-test-api
   ```

## Running the Server

Start the server with:

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### POST /hello_world

Accepts a JSON payload with a name and returns a greeting message.

**Request:**

```bash
curl -X POST http://localhost:8080/hello_world \
  -H "Content-Type: application/json" \
  -d '{"name": "World"}'
```

**Request Body:**

```json
{
  "name": "World"
}
```

**Response:**

```json
{
  "message": "Hello, World!"
}
```

## Example Usage

```bash
# Basic hello world
curl -X POST http://localhost:8080/hello_world \
  -H "Content-Type: application/json" \
  -d '{"name": "World"}'

# Custom name
curl -X POST http://localhost:8080/hello_world \
  -H "Content-Type: application/json" \
  -d '{"name": "Aaron"}'
```

## Building

To build a binary:

```bash
go build -o api-server
./api-server
```
