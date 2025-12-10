package books

import (
	"context"
	"errors"
	"testing"
)

func TestDefaultBooksService_GetMetrics_Success(t *testing.T) {
	mock := NewMockBooksProvider()
	svc := NewBooksService(mock)

	metrics, err := svc.GetMetrics(context.Background(), "Alan Donovan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metrics.MeanUnitsSold != 5000 {
		t.Fatalf("unexpected mean units sold: %d", metrics.MeanUnitsSold)
	}
	if metrics.CheapestBook != "The Go Programming Language" {
		t.Fatalf("unexpected cheapest book: %q", metrics.CheapestBook)
	}
	if metrics.BooksWrittenByAuthor != 1 {
		t.Fatalf("unexpected books written by author: %d", metrics.BooksWrittenByAuthor)
	}
}

func TestDefaultBooksService_GetMetrics_Empty(t *testing.T) {
	mock := NewMockBooksProviderWith([]Book{}, nil)
	svc := NewBooksService(mock)

	_, err := svc.GetMetrics(context.Background(), "Any")
	if err == nil {
		t.Fatalf("expected error for empty books, got nil")
	}
}

func TestDefaultBooksService_GetMetrics_ProviderError(t *testing.T) {
	wantErr := errors.New("provider failure")
	mock := NewMockBooksProviderWith(nil, wantErr)
	svc := NewBooksService(mock)

	_, err := svc.GetMetrics(context.Background(), "Any")
	if err == nil {
		t.Fatalf("expected provider error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected wrapped error to contain provider error; got: %v", err)
	}
}
