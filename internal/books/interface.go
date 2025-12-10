package books

import "context"

type BooksProvider interface {
	GetBooks(ctx context.Context) ([]Book, error)
}
