package handlers

import (
	"context"
	"net/http"
	"slices"

	"educabot.com/bookshop/internal/models"
	"educabot.com/bookshop/internal/providers"
	"github.com/gin-gonic/gin"
)

type GetMetricsRequest struct {
	Author string `form:"author"`
}

func NewGetMetrics(booksProvider providers.BooksProvider) GetMetrics {
	return GetMetrics{booksProvider}
}

type GetMetrics struct {
	booksProvider providers.BooksProvider
}

func (h GetMetrics) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query GetMetricsRequest
		ctx.ShouldBindQuery(&query) // this returns error

		books, err := h.booksProvider.GetBooks(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusBadGateway, map[string]interface{}{"error": err.Error()})
			return
		}

		meanUnitsSold := meanUnitsSold(ctx.Request.Context(), books)
		cheapestBook := cheapestBook(ctx.Request.Context(), books).Name
		booksWrittenByAuthor := booksWrittenByAuthor(ctx.Request.Context(), books, query.Author)

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"mean_units_sold":         meanUnitsSold,
			"cheapest_book":           cheapestBook,
			"books_written_by_author": booksWrittenByAuthor,
		})
	}
}

func meanUnitsSold(_ context.Context, books []models.Book) uint {
	var sum uint
	for _, book := range books {
		sum += book.UnitsSold
	}
	return sum / uint(len(books))
}

func cheapestBook(_ context.Context, books []models.Book) models.Book {
	return slices.MinFunc(books, func(a, b models.Book) int {
		return int(a.Price - b.Price)
	})
}

func booksWrittenByAuthor(_ context.Context, books []models.Book, author string) uint {
	var count uint
	for _, book := range books {
		if book.Author == author {
			count++
		}
	}
	return count
}
