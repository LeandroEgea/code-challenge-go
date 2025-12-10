package books

import "context"

// MockBooksService is a simple configurable mock of the BooksService.
// It returns the preset Metrics and error when GetMetrics is called.
type MockBooksService struct {
	Metrics Metrics
	Err     error
}

// NewMockBooksService creates a mock service that will return the provided
// metrics and error when invoked.
func NewMockBooksService(m Metrics, err error) *MockBooksService {
	return &MockBooksService{Metrics: m, Err: err}
}

// GetMetrics implements BooksService.
func (m *MockBooksService) GetMetrics(_ context.Context, _ string) (Metrics, error) {
	return m.Metrics, m.Err
}
