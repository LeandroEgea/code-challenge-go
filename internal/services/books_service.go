package services

import (
	"context"
	"fmt"
	"slices"

	"educabot.com/bookshop/internal/models"
	"educabot.com/bookshop/internal/providers"
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
	provider providers.BooksProvider
}

func NewBooksService(p providers.BooksProvider) *DefaultBooksService {
	return &DefaultBooksService{provider: p}
}

func (s *DefaultBooksService) GetMetrics(ctx context.Context, author string) (Metrics, error) {
	books, err := s.provider.GetBooks(ctx)
	if err != nil {
		return Metrics{}, fmt.Errorf("fetching books: %w", err)
	}

	if len(books) == 0 {
		return Metrics{}, fmt.Errorf("no books available")
	}

	// mean units sold
	var sum uint
	for _, b := range books {
		sum += b.UnitsSold
	}
	mean := sum / uint(len(books))

	// cheapest
	cheapest := slices.MinFunc(books, func(a, b models.Book) int {
		return int(a.Price - b.Price)
	})

	// count by author
	var count uint
	for _, b := range books {
		if b.Author == author {
			count++
		}
	}

	return Metrics{
		MeanUnitsSold:        mean,
		CheapestBook:         cheapest.Name,
		BooksWrittenByAuthor: count,
	}, nil
}

func meanUnitsSold(books []models.Book) uint {
	var sum uint
	for _, book := range books {
		sum += book.UnitsSold
	}
	return sum / uint(len(books))
}

func cheapestBook(books []models.Book) models.Book {
	return slices.MinFunc(books, func(a, b models.Book) int {
		return int(a.Price - b.Price)
	})
}

func booksWrittenByAuthor(books []models.Book, author string) uint {
	var count uint
	for _, book := range books {
		if book.Author == author {
			count++
		}
	}
	return count
}
