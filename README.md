# Bookshop (Go)

A small Go service that fetches books from a remote provider and computes simple metrics (mean units sold, cheapest book, and books written by an author). The project is organized with a clean separation between provider, service, and HTTP handler layers and includes unit and integration tests.

Project structure (high level):
- `cmd/main` — application entrypoint and HTTP handlers
- `internal/books` — domain models, provider (HTTP), service (business logic), and mocks for tests
- `static` — simple UI assets (not required for the service)

Running locally:

1. Ensure Go is installed (1.21+ preferred).
2. From the project root run the tests:

```powershell
go test ./... 
```

3. To run the server locally:

```powershell
go run ./cmd/main
# then open: http://localhost:3000/?author=J.R.R.%20Tolkien
```

Tests:
- Unit tests cover the HTTP provider, the service logic, and the HTTP handlers.
- An integration test exercises the service against a mocked “real” API response to validate end-to-end behavior.

Next steps (short):
- Implement a caching layer in the service: fetch from the `HTTPBooksProvider` once and cache the result for subsequent calls (TTL or in-memory cache) to reduce remote calls.
- Add metrics and structured logging: capture request counts, latency, and provider errors; emit structured logs (e.g., JSON) to make debugging and observability easier.
