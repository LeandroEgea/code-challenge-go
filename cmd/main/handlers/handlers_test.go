package handlers

import (
	"encoding/json"
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

	assert.Equal(t, 11000, int(resBody["mean_units_sold"].(float64)))
	assert.Equal(t, "The Go Programming Language", resBody["cheapest_book"])
	assert.Equal(t, 1, int(resBody["books_written_by_author"].(float64)))
}
