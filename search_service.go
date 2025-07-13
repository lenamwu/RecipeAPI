package main

import (
	"math"
	"sort"
	"strings"
)

// SearchService handles recipe search operations
type SearchService struct {
	recipes map[int]*RecipeWithStats
}

// NewSearchService creates a new SearchService instance
func NewSearchService(recipes map[int]*RecipeWithStats) *SearchService {
	return &SearchService{
		recipes: recipes,
	}
}

// SearchRecipes searches for recipes matching the given query
// Only searches in recipe titles/names and returns top 20 results
func (ss *SearchService) SearchRecipes(query string) []Recipe {
	if query == "" {
		return []Recipe{}
	}

	queryLower := strings.ToLower(query)
	var matchingRecipes []Recipe

	// Search through recipes
	for _, recipeStats := range ss.recipes {
		// Only include recipes with reviews and good ratings (>= 4.0)
		if recipeStats.ReviewCount == 0 || recipeStats.AvgRating < 4.0 {
			continue
		}

		// Check if query matches ONLY the recipe title/name
		if strings.Contains(strings.ToLower(recipeStats.Name), queryLower) {
			matchingRecipes = append(matchingRecipes, recipeStats.Recipe)
		}
	}

	// Sort by weighted score that considers both rating and review count
	sort.Slice(matchingRecipes, func(i, j int) bool {
		scoreI := ss.calculateWeightedScore(matchingRecipes[i].AvgRating, matchingRecipes[i].NReviews)
		scoreJ := ss.calculateWeightedScore(matchingRecipes[j].AvgRating, matchingRecipes[j].NReviews)
		return scoreI > scoreJ
	})

	// Limit to top 20 results
	if len(matchingRecipes) > 20 {
		matchingRecipes = matchingRecipes[:20]
	}

	return matchingRecipes
}

// calculateWeightedScore computes a weighted score that balances rating and review count
// This ensures recipes with many reviews don't get unfairly penalized by slightly lower ratings
func (ss *SearchService) calculateWeightedScore(avgRating float64, nReviews int) float64 {
	// Base score from rating (0-5 scale)
	ratingScore := avgRating

	// Confidence factor based on number of reviews
	// Uses a logarithmic scale to prevent recipes with thousands of reviews from dominating
	// but still gives significant weight to review count
	confidenceFactor := math.Log10(float64(nReviews) + 1)

	// Weighted score: rating * confidence factor
	// A recipe with 2000 reviews at 4.8 rating will score higher than 3 reviews at 5.0
	// Example: 4.8 * log10(2001) ≈ 4.8 * 3.3 ≈ 15.8
	// vs: 5.0 * log10(4) ≈ 5.0 * 0.6 ≈ 3.0
	return ratingScore * confidenceFactor
}

// GetStats returns statistics about the recipe collection
func (ss *SearchService) GetStats() (int, int) {
	totalRecipes := len(ss.recipes)
	goodRecipes := 0

	for _, recipe := range ss.recipes {
		if recipe.ReviewCount > 0 && recipe.AvgRating >= 4.0 {
			goodRecipes++
		}
	}

	return totalRecipes, goodRecipes
}
