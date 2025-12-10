package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"educabot.com/bookshop/internal/models"
)

type HTTPBooksProvider struct {
	endpoint string
	client   *http.Client
}

// NewHTTPBooksProvider creates a provider that fetches books from the given endpoint.
func NewHTTPBooksProvider(endpoint string) *HTTPBooksProvider {
	return &HTTPBooksProvider{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetBooks fetches books from the remote endpoint, validates the payload and returns parsed books or an error.
func (p *HTTPBooksProvider) GetBooks(ctx context.Context) ([]models.Book, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// read body for debugging but cap the size
		limited, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(limited))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	// Temporary struct to be defensive about types returned by the remote API
	var raw []struct {
		ID        json.Number `json:"id"`
		Name      string      `json:"name"`
		Author    string      `json:"author"`
		UnitsSold json.Number `json:"units_sold"`
		Price     json.Number `json:"price"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("invalid json: %w", err)
	}

	books := make([]models.Book, 0, len(raw))
	for i, r := range raw {
		// parse numeric fields defensively
		id64, err := r.ID.Int64()
		if err != nil {
			return nil, fmt.Errorf("invalid id at index %d: %w", i, err)
		}
		units64, err := r.UnitsSold.Int64()
		if err != nil {
			return nil, fmt.Errorf("invalid units_sold at index %d: %w", i, err)
		}
		price64, err := r.Price.Int64()
		if err != nil {
			return nil, fmt.Errorf("invalid price at index %d: %w", i, err)
		}

		if r.Name == "" || r.Author == "" {
			return nil, fmt.Errorf("invalid payload at index %d: missing name or author", i)
		}
		if units64 < 0 || price64 < 0 {
			return nil, fmt.Errorf("invalid numeric values at index %d", i)
		}

		books = append(books, models.Book{
			ID:        uint(id64),
			Name:      r.Name,
			Author:    r.Author,
			UnitsSold: uint(units64),
			Price:     uint(price64),
		})
	}

	return books, nil
}
