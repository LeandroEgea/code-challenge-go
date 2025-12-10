package books

import "context"

type BooksService interface {
	GetMetrics(ctx context.Context, author string) (Metrics, error)
}

type BooksProvider interface {
	GetBooks(ctx context.Context) ([]Book, error)
}
