package main

import (
	"encoding/json"
	"io"
	"log"
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
		"status":              "healthy",
		"total_recipes":       totalRecipes,
		"recipes_with_rating": recipesWithRating,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ImageProxyHandler serves images via proxy to bypass CORS issues
func (h *Handlers) ImageProxyHandler(w http.ResponseWriter, r *http.Request) {
	imgURL := r.URL.Query().Get("url")
	if imgURL == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(imgURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch image", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(http.StatusOK)

	_, copyErr := io.Copy(w, resp.Body)
	if copyErr != nil {
		log.Println("Error copying image data:", copyErr)
	}
}
