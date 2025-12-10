package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"educabot.com/bookshop/internal/books"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetMetrics_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := books.NewBooksService(books.NewMockBooksProvider())
	handler := NewGetMetrics(service)

	r := gin.Default()
	r.GET("/", handler.Handle())

	req := httptest.NewRequest(http.MethodGet, "/?author=Alan+Donovan", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var resBody map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(t, 5000, int(resBody["mean_units_sold"].(float64)))
	assert.Equal(t, "The Go Programming Language", resBody["cheapest_book"])
	assert.Equal(t, 1, int(resBody["books_written_by_author"].(float64)))
}

func TestGetMetrics_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := books.NewMockBooksService(books.Metrics{}, errors.New("service error"))
	handler := NewGetMetrics(mockService)

	r := gin.Default()
	r.GET("/", handler.Handle())

	req := httptest.NewRequest(http.MethodGet, "/?author=test", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadGateway, res.Code)

	var resBody map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.NotNil(t, resBody["error"])
}

func TestGetMetrics_NoAuthorFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := books.NewBooksService(books.NewMockBooksProvider())
	handler := NewGetMetrics(service)

	r := gin.Default()
	r.GET("/", handler.Handle())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var resBody map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, 11000, int(resBody["mean_units_sold"].(float64)))
	assert.Equal(t, 0, int(resBody["books_written_by_author"].(float64)))
}
