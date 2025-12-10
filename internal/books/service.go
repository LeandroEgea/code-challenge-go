package books

import (
	"context"
	"fmt"
	"slices"
)

type Metrics struct {
	MeanUnitsSold        uint   `json:"mean_units_sold"`
	CheapestBook         string `json:"cheapest_book"`
	BooksWrittenByAuthor uint   `json:"books_written_by_author"`
}

type BooksService interface {
	GetMetrics(ctx context.Context, author string) (Metrics, error)
}

type DefaultBooksService struct {
	provider BooksProvider
}

func NewBooksService(p BooksProvider) *DefaultBooksService {
	return &DefaultBooksService{provider: p}
}

func (s *DefaultBooksService) GetMetrics(ctx context.Context, author string) (Metrics, error) {
	books, err := s.provider.GetBooks(ctx)
	if err != nil {
		return Metrics{}, fmt.Errorf("fetching books: %w", err)
	}

	// TODO: Consider implementing a caching layer here to avoid repeated calls to the provider.
	// Cache could be time-based (TTL) or invalidated on demand via a separate method.

	if len(books) == 0 {
		return Metrics{}, fmt.Errorf("no books available")
	}

	mean := meanUnitsSold(books)
	cheapest := cheapestBook(books)
	count := booksWrittenByAuthor(books, author)

	return Metrics{
		MeanUnitsSold:        mean,
		CheapestBook:         cheapest.Name,
		BooksWrittenByAuthor: count,
	}, nil
}

func meanUnitsSold(books []Book) uint {
	var sum uint
	for _, book := range books {
		sum += book.UnitsSold
	}
	return sum / uint(len(books))
}

func cheapestBook(books []Book) Book {
	return slices.MinFunc(books, func(a, b Book) int {
		return int(a.Price - b.Price)
	})
}

func booksWrittenByAuthor(books []Book, author string) uint {
	var count uint
	for _, book := range books {
		if book.Author == author {
			count++
		}
	}
	return count
}
