package main

import (
	"encoding/json"
	"net/http"
)

// Handlers contains all HTTP request handlers
type Handlers struct {
	searchService *SearchService
	dataLoader    *DataLoader
}

// NewHandlers creates a new Handlers instance
func NewHandlers(searchService *SearchService, dataLoader *DataLoader) *Handlers {
	return &Handlers{
		searchService: searchService,
		dataLoader:    dataLoader,
	}
}

// SearchRecipesHandler handles recipe search requests
func (h *Handlers) SearchRecipesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	matchingRecipes := h.searchService.SearchRecipes(query)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"recipes": matchingRecipes,
		"count":   len(matchingRecipes),
		"query":   query,
	}
	json.NewEncoder(w).Encode(response)
}

// HealthCheckHandler handles health check requests
func (h *Handlers) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	totalRecipes, recipesWithRating := h.searchService.GetStats()

	response := map[string]interface{}{
		"status":               "healthy",
		"total_recipes":        totalRecipes,
		"recipes_with_rating":  recipesWithRating,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
