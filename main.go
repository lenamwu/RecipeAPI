package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Initialize data loader
	dataLoader := NewDataLoader()

	// Load data from CSV files
	if err := dataLoader.LoadData(); err != nil {
		log.Fatal("Failed to load data:", err)
	}

	// Initialize search service
	searchService := NewSearchService(dataLoader.GetRecipes())

	// Initialize handlers
	handlers := NewHandlers(searchService, dataLoader)

	// Setup HTTP routes
	http.HandleFunc("/recipes", handlers.SearchRecipesHandler)
	http.HandleFunc("/health", handlers.HealthCheckHandler)
	http.HandleFunc("/img", handlers.ImageProxyHandler)

	// Get the port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running at http://localhost:%s\n", port)
	log.Printf("Try: http://localhost:%s/recipes?query=apple\n", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
