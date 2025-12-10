# Book Metrics API

A domain-driven design Go application that fetches book data from a remote API and computes metrics such as mean units sold, cheapest book, and books written by a specific author.

## Architecture

The application follows a domain-driven design with clear separation of concerns:

```
cmd/main/
├── main.go              # Application entry point
└── handlers/
    ├── handlers.go      # HTTP request handlers
    └── handlers_test.go # Handler unit tests

internal/books/
├── model.go             # Domain models (Book, Metrics)
├── interface.go         # Provider and Service interfaces
├── repository.go        # HTTPBooksProvider (external data fetch)
├── repository_mock.go   # Mock provider for testing
├── repository_http_test.go # HTTP provider unit tests
├── service.go           # DefaultBooksService (business logic)
├── service_mock.go      # Mock service for handler testing
└── service_test.go      # Service unit tests

static/
├── index.html           # Static HTML
└── css/style.css        # Styling
```

### Key Components

- **Provider** (`BooksProvider`): Fetches books from external sources (HTTP or mock).
- **Service** (`BooksService`): Computes metrics from a list of books.
- **Handler**: HTTP request handler that orchestrates provider and service.
- **Models**: `Book` and `Metrics` domain entities.

## Running Locally

### Prerequisites
- Go 1.21 or later
- `curl` or Postman for testing

### 1. Run All Tests
```powershell
go test ./...
```

This runs:
- Handler tests (3 cases: success, service error, no author filter)
- Service tests (3 cases: success, empty result, provider error)
- Repository HTTP tests (3 cases: success, non-200 status, invalid JSON)

### 2. Start the Application
```powershell
go run ./cmd/main
```

The server starts on `http://localhost:3000`

### 3. Test the API Endpoint

**Query with author filter:**
```powershell
curl "http://localhost:3000/?author=Alan%20Donovan"
```

**Query without author filter:**
```powershell
curl "http://localhost:3000/"
```

**Query with J.R.R. Tolkien (if available):**
```powershell
curl "http://localhost:3000/?author=J.R.R.%20Tolkien"
```

### Expected Response
```json
{
  "mean_units_sold": 11000,
  "cheapest_book": "The Go Programming Language",
  "books_written_by_author": 1
}
```

### 4. View Available Books
Check what books and authors are in the remote endpoint:
```powershell
curl "https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books"
```

## Next Steps (Priority Order)

### ✅ Fixed
- **Cheapest book calculation**: The `cheapestBook()` function now correctly uses explicit comparison logic instead of subtraction, ensuring it returns the book with the minimum price.

### 🔧 TODO: High Priority

#### 1. Fix Mean Units Sold Edge Case
**Issue**: The mean calculation uses integer division, which truncates fractional parts.

**Current code:**
```go
func meanUnitsSold(books []Book) uint {
    var sum uint
    for _, book := range books {
        sum += book.UnitsSold
    }
    return sum / uint(len(books))  // Integer division truncates
}
```

**Improvement**: Return `float64` or use rounding:
```go
return (sum + uint(len(books))/2) / uint(len(books))  // Round to nearest
```

#### 2. Add Integration Tests with Real Endpoint
Create `internal/books/integration_test.go` that:
- Fetches real data from the endpoint
- Validates response structure
- Tests edge cases (empty results, malformed data)
- Runs against the live API

Example:
```go
func TestHTTPBooksProvider_RealEndpoint(t *testing.T) {
    p := NewHTTPBooksProvider("https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books")
    books, err := p.GetBooks(context.Background())
    
    if err != nil {
        t.Fatalf("failed to fetch from real endpoint: %v", err)
    }
    if len(books) == 0 {
        t.Fatal("expected books from real endpoint")
    }
}
```

#### 3. Implement Caching Layer
Add a cache wrapper in `internal/books/cache.go`:

```go
type CachedBooksService struct {
    underlying BooksService
    cache      []Book
    cacheTTL   time.Duration
    lastFetch  time.Time
    mu         sync.RWMutex
}

func (c *CachedBooksService) GetMetrics(ctx context.Context, author string) (Metrics, error) {
    c.mu.RLock()
    if time.Since(c.lastFetch) < c.cacheTTL && c.cache != nil {
        c.mu.RUnlock()
        // Use cached data
        return computeMetrics(c.cache, author), nil
    }
    c.mu.RUnlock()
    
    // Fetch fresh data and update cache
    metrics, err := c.underlying.GetMetrics(ctx, author)
    if err == nil {
        c.mu.Lock()
        c.cache = books  // Store in cache
        c.lastFetch = time.Now()
        c.mu.Unlock()
    }
    return metrics, err
}
```

Wire it in `cmd/main/main.go`:
```go
remote := books.NewHTTPBooksProvider(endpoint)
cached := books.NewCachedBooksService(remote, 5*time.Minute)
service := books.NewBooksService(cached)
```

### 🎯 Medium Priority

#### 4. Configuration Management
Add env var support in `cmd/main/main.go`:
```go
endpoint := os.Getenv("BOOKS_API_URL")
if endpoint == "" {
    endpoint = "https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books"
}
```

#### 5. Structured Logging
Add logging for debugging and production monitoring:
```go
log.Printf("Fetching books from %s", endpoint)
log.Printf("Service computed metrics: %+v", metrics)
```

#### 6. Error Type Mapping
Distinguish HTTP error codes based on error type:
- `504 Gateway Timeout` for timeouts
- `503 Service Unavailable` for provider errors
- `400 Bad Request` for parsing errors

### 📚 Lower Priority

#### 7. Enhanced Index Page
Update `static/index.html` with:
- Form to query the metrics endpoint
- Display results in a formatted table
- Show available authors

#### 8. Documentation
Add package-level godoc comments:
```go
// Package books provides domain models, data access, and business logic
// for book management and metrics computation.
package books
```

## Testing Strategy

### Unit Tests
All services, handlers, and providers have unit tests using mocks.

```powershell
go test -v ./...
```

### Integration Tests (TODO)
Test against the real API endpoint to catch breaking changes.

```powershell
go test -v -tags=integration ./...
```

## Future Enhancements

- [ ] Graceful shutdown (SIGINT/SIGTERM handling)
- [ ] Middleware for logging, tracing, and panic recovery
- [ ] Rate limiting on the HTTP provider
- [ ] Metrics export (Prometheus)
- [ ] Dependency injection framework (if complexity grows)

## Troubleshooting

### Program returns "A Wizard of Earthsea" for all queries
**Status**: ✅ Fixed in the latest version.

**If still occurring**: The `cheapestBook()` function's comparison logic was reversed. Ensure you have the latest code:
```go
func cheapestBook(books []Book) Book {
    return slices.MinFunc(books, func(a, b Book) int {
        if a.Price < b.Price {
            return -1
        }
        if a.Price > b.Price {
            return 1
        }
        return 0
    })
}
```

### Mean units sold is incorrect
The current implementation truncates fractional parts due to integer division. This will be addressed in TODO #1.

### Server won't start
- Check if port 3000 is already in use: `netstat -ano | findstr :3000`
- Ensure the remote endpoint is reachable: `curl https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books`

## License
MIT