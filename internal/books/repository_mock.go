package books

import "context"

// MockBooksProvider is a configurable mock implementation of the BooksProvider.
type MockBooksProvider struct {
	Books []Book
	Err   error
}

// NewMockBooksProvider returns a mock populated with sample data (convenience).
func NewMockBooksProvider() *MockBooksProvider {
	return &MockBooksProvider{
		Books: []Book{
			{ID: 1, Name: "The Go Programming Language", Author: "Alan Donovan", UnitsSold: 5000, Price: 40},
			{ID: 2, Name: "Clean Code", Author: "Robert C. Martin", UnitsSold: 15000, Price: 50},
			{ID: 3, Name: "The Pragmatic Programmer", Author: "Andrew Hunt", UnitsSold: 13000, Price: 45},
		},
		Err: nil,
	}
}

// NewMockBooksProviderWith returns a mock configured with the given books and error.
func NewMockBooksProviderWith(books []Book, err error) *MockBooksProvider {
	return &MockBooksProvider{Books: books, Err: err}
}

func (m *MockBooksProvider) GetBooks(_ context.Context) ([]Book, error) {
	return m.Books, m.Err
}
