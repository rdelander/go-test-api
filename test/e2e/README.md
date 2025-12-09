# E2E Tests

End-to-end tests that verify the complete application stack running in Docker.

## Running Locally

Run all E2E tests:

```bash
make test-e2e
```

This will:

1. Build the Docker image
2. Start PostgreSQL container
3. Start API container (with auto-migrations)
4. Run Go E2E tests against the live API
5. Clean up containers automatically

## Running in CI

The CI workflow (`test-e2e-ci` target) assumes services are already running and just executes the tests.

## Test Structure

Tests are in `test/e2e/` with the build tag `e2e`:

```go
//go:build e2e
// +build e2e
```

This ensures they only run when explicitly requested.

## Writing E2E Tests

Tests use standard Go testing with real HTTP requests:

```go
func TestExample(t *testing.T) {
    // API is already running at baseURL
    resp, err := http.Get(baseURL + "/endpoint")
    // ... assertions
}
```

The `waitForAPI()` helper ensures the API is ready before tests run.

## Environment Variables

- `API_BASE_URL`: Base URL for the API (default: `http://localhost:8080`)

## What Gets Tested

- ✅ Real Docker image (same one deployed to production)
- ✅ Database migrations run on startup
- ✅ HTTP endpoints with real PostgreSQL
- ✅ Request/response format validation
- ✅ Error cases and validation

## Tests Included

- `TestHealthEndpoint` - Verifies health check works
- `TestRegisterUser` - Creates a new user via API
- `TestRegisterUserDuplicate` - Validates duplicate email handling
