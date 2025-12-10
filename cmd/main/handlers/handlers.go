package handlers

import (
	"net/http"

	"educabot.com/bookshop/internal/books"
	"github.com/gin-gonic/gin"
)

type GetMetricsRequest struct {
	Author string `form:"author"`
}

func NewGetMetrics(booksService books.BooksService) GetMetrics {
	return GetMetrics{booksService}
}

type GetMetrics struct {
	booksService books.BooksService
}

func (h GetMetrics) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query GetMetricsRequest
		if err := ctx.ShouldBindQuery(&query); err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid query parameters"})
			return
		}

		metrics, err := h.booksService.GetMetrics(ctx.Request.Context(), query.Author)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, map[string]interface{}{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, metrics)
	}
}
