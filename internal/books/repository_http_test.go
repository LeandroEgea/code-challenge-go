package books_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"educabot.com/bookshop/internal/books"
)

func TestHTTPBooksProvider_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
            {"id":1,"name":"The Go Programming Language","author":"Alan Donovan","units_sold":5000,"price":40},
            {"id":2,"name":"Clean Code","author":"Robert C. Martin","units_sold":15000,"price":50}
        ]`))
	}))
	defer srv.Close()

	p := books.NewHTTPBooksProvider(srv.URL)
	got, err := p.GetBooks(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 books, got %d", len(got))
	}
	if got[0].Name != "The Go Programming Language" {
		t.Fatalf("unexpected first book name: %q", got[0].Name)
	}
}

func TestHTTPBooksProvider_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer srv.Close()

	p := books.NewHTTPBooksProvider(srv.URL)
	_, err := p.GetBooks(context.Background())
	if err == nil {
		t.Fatalf("expected error for non-200 response, got nil")
	}
}

func TestHTTPBooksProvider_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not a json"))
	}))
	defer srv.Close()

	p := books.NewHTTPBooksProvider(srv.URL)
	_, err := p.GetBooks(context.Background())
	if err == nil {
		t.Fatalf("expected error for invalid json, got nil")
	}
}
