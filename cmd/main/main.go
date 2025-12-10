package main

import (
	"fmt"

	handlers "educabot.com/bookshop/cmd/main/handlers"
	"educabot.com/bookshop/internal/providers"
	"educabot.com/bookshop/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.SetTrustedProxies(nil)

	// Use the remote HTTP provider as the primary source of books.
	remote := providers.NewHTTPBooksProvider("https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books")
	service := services.NewBooksService(remote)
	metricsHandler := handlers.NewGetMetrics(service)
	router.GET("/", metricsHandler.Handle())
	router.Run(":3000")
	fmt.Println("Starting server on :3000")
}
