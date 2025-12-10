package books

import (
	"context"
	"fmt"
	"slices"
)

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

	// Filter books by author if specified
	var booksToAnalyze []Book
	if author != "" {
		for _, b := range books {
			if b.Author == author {
				booksToAnalyze = append(booksToAnalyze, b)
			}
		}
		// If author filter specified but no books found, return metrics with zeros
		if len(booksToAnalyze) == 0 {
			return Metrics{
				MeanUnitsSold:        0,
				CheapestBook:         "",
				BooksWrittenByAuthor: 0,
			}, nil
		}
	} else {
		booksToAnalyze = books
	}

	mean := meanUnitsSold(booksToAnalyze)
	cheapest := cheapestBook(booksToAnalyze)
	count := booksWrittenByAuthor(booksToAnalyze, author)

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
		if a.Price < b.Price {
			return -1
		}
		if a.Price > b.Price {
			return 1
		}
		return 0
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
