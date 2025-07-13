package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize data loader
	dataLoader := NewDataLoader()

	// Load data from CSV files
	if err := dataLoader.LoadData(); err != nil {
		log.Fatal("Failed to load data:", err)
	}

	// Calculate recipe statistics
	dataLoader.CalculateRecipeStats()

	// Initialize search service
	searchService := NewSearchService(dataLoader.GetRecipes())

	// Initialize handlers
	handlers := NewHandlers(searchService, dataLoader)

	// Setup HTTP routes
	http.HandleFunc("/recipes", handlers.SearchRecipesHandler)
	http.HandleFunc("/health", handlers.HealthCheckHandler)

	log.Println("Server running at http://localhost:8080")
	log.Println("Try: http://localhost:8080/recipes?query=chicken")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
