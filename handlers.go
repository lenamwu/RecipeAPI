package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
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

    // Rewrite img_src to use your proxy
    for i := range matchingRecipes {
        originalURL := matchingRecipes[i].ImageSrc
        if originalURL != "" {
            encodedURL := url.QueryEscape(originalURL)
            matchingRecipes[i].ImageSrc = "https://recipeapi-1-b1hi.onrender.com/img?url=" + encodedURL
        }
    }

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

func (h *Handlers) ImageProxyHandler(w http.ResponseWriter, r *http.Request) {
	imgURL := r.URL.Query().Get("url")
	if imgURL == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	// Create new HTTP request
	req, err := http.NewRequest("GET", imgURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusBadGateway)
		return
	}

	// Set User-Agent header (pretend to be a browser)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RecipeAPI/1.0)")

	// Perform request with default client
	resp, err := http.DefaultClient.Do(req)
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
