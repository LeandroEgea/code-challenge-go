package books_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"educabot.com/bookshop/internal/books"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_RealAPIResponse(t *testing.T) {
	// Mock server with real API response
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"id":1,"name":"The Fellowship of the Ring","author":"J.R.R. Tolkien","units_sold":50000000,"price":20},
			{"id":2,"name":"The Two Towers","author":"J.R.R. Tolkien","units_sold":30000000,"price":20},
			{"id":3,"name":"The Return of the King","author":"J.R.R. Tolkien","units_sold":50000000,"price":20},
			{"id":4,"name":"The Lion, the Witch and the Wardrobe","author":"C.S. Lewis","units_sold":85000000,"price":15},
			{"id":5,"name":"A Wizard of Earthsea","author":"Ursula K. Le Guin","units_sold":1000000,"price":10},
			{"id":6,"name":"The Hobbit","author":"J.R.R. Tolkien","units_sold":140000000,"price":25}
		]`))
	}))
	defer srv.Close()

	// Test provider
	provider := books.NewHTTPBooksProvider(srv.URL)
	booksData, err := provider.GetBooks(context.Background())
	if err != nil {
		t.Fatalf("provider error: %v", err)
	}

	t.Logf("Fetched %d books", len(booksData))
	assert.Equal(t, 6, len(booksData), "expected 6 books from mocked endpoint")
	for _, b := range booksData {
		t.Logf("Book: %s (Author: %s, Price: %d, Units: %d)", b.Name, b.Author, b.Price, b.UnitsSold)
	}

	// Test service with real data
	service := books.NewBooksService(provider)

	// Test 1: Tolkien metrics
	t.Run("Tolkien_Metrics", func(t *testing.T) {
		metrics, err := service.GetMetrics(context.Background(), "J.R.R. Tolkien")
		assert.NoError(t, err, "service should not error for Tolkien")
		assert.Equal(t, uint(4), metrics.BooksWrittenByAuthor, "Tolkien should have 4 books")
		assert.Equal(t, "The Fellowship of the Ring", metrics.CheapestBook, "cheapest Tolkien book should be Fellowship (price 20)")
		// Mean of [50000000, 30000000, 50000000, 140000000] = 67500000
		expectedMean := uint(67500000)
		assert.Equal(t, expectedMean, metrics.MeanUnitsSold, "mean units sold for Tolkien should be 67500000")
		t.Logf("Tolkien metrics: mean=%d, cheapest=%s, count=%d",
			metrics.MeanUnitsSold, metrics.CheapestBook, metrics.BooksWrittenByAuthor)
	})

	// Test 2: C.S. Lewis metrics
	t.Run("CSLewis_Metrics", func(t *testing.T) {
		metrics, err := service.GetMetrics(context.Background(), "C.S. Lewis")
		assert.NoError(t, err, "service should not error for C.S. Lewis")
		assert.Equal(t, uint(1), metrics.BooksWrittenByAuthor, "C.S. Lewis should have 1 book")
		assert.Equal(t, "The Lion, the Witch and the Wardrobe", metrics.CheapestBook, "cheapest book overall is A Wizard of Earthsea (price 10)")
		t.Logf("C.S. Lewis metrics: mean=%d, cheapest=%s, count=%d",
			metrics.MeanUnitsSold, metrics.CheapestBook, metrics.BooksWrittenByAuthor)
	})

	// Test 3: Ursula K. Le Guin metrics
	t.Run("LeGuin_Metrics", func(t *testing.T) {
		metrics, err := service.GetMetrics(context.Background(), "Ursula K. Le Guin")
		assert.NoError(t, err, "service should not error for Ursula K. Le Guin")
		assert.Equal(t, uint(1), metrics.BooksWrittenByAuthor, "Le Guin should have 1 book")
		assert.Equal(t, "A Wizard of Earthsea", metrics.CheapestBook, "cheapest book overall is A Wizard of Earthsea (price 10)")
		t.Logf("Le Guin metrics: mean=%d, cheapest=%s, count=%d",
			metrics.MeanUnitsSold, metrics.CheapestBook, metrics.BooksWrittenByAuthor)
	})

	// Test 4: No author filter
	t.Run("NoAuthor_Metrics", func(t *testing.T) {
		metrics, err := service.GetMetrics(context.Background(), "")
		assert.NoError(t, err, "service should not error for empty author filter")
		assert.Equal(t, uint(0), metrics.BooksWrittenByAuthor, "empty author should match 0 books")
		assert.Equal(t, "A Wizard of Earthsea", metrics.CheapestBook, "cheapest book is A Wizard of Earthsea (price 10)")
		// Mean of all: [50000000, 30000000, 50000000, 85000000, 1000000, 140000000] = 59333333.33 (truncated to 59333333)
		expectedMean := uint(59333333)
		assert.Equal(t, expectedMean, metrics.MeanUnitsSold, "mean units sold across all books should be 59333333")
		t.Logf("All books metrics: mean=%d, cheapest=%s, count=%d",
			metrics.MeanUnitsSold, metrics.CheapestBook, metrics.BooksWrittenByAuthor)
	})

	// Test 5: Non-existent author
	t.Run("NonExistent_Author", func(t *testing.T) {
		metrics, err := service.GetMetrics(context.Background(), "Non Existent")
		assert.NoError(t, err, "service should not error for non-existent author")
		assert.Equal(t, uint(0), metrics.BooksWrittenByAuthor, "non-existent author should match 0 books")
		assert.Equal(t, "A Wizard of Earthsea", metrics.CheapestBook, "cheapest book is A Wizard of Earthsea (price 10)")
		t.Logf("Non-existent author metrics: mean=%d, cheapest=%s, count=%d",
			metrics.MeanUnitsSold, metrics.CheapestBook, metrics.BooksWrittenByAuthor)
	})
}
