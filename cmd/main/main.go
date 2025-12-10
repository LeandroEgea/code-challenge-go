package main

import (
	"fmt"
	"os"

	handlers "educabot.com/bookshop/cmd/main/handlers"
	"educabot.com/bookshop/internal/books"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	if err := router.SetTrustedProxies(nil); err != nil {
		fmt.Fprintf(os.Stderr, "error setting trusted proxies: %v\n", err)
		os.Exit(1)
	}

	// Use the remote HTTP provider as the primary source of books.
	remote := books.NewHTTPBooksProvider("https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books")
	service := books.NewBooksService(remote)
	metricsHandler := handlers.NewGetMetrics(service)
	router.GET("/", metricsHandler.Handle())
	fmt.Println("Starting server on :3000")
	if err := router.Run(":3000"); err != nil {
		fmt.Fprintf(os.Stderr, "error running router: %v\n", err)
		os.Exit(1)
	}
}
