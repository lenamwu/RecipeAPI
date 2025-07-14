package main

import (
	"sort"
	"strings"
	"unicode"
)

// SearchService handles recipe search operations
type SearchService struct {
	recipes map[int]*Recipe
}

// NewSearchService creates a new SearchService instance
func NewSearchService(recipes map[int]*Recipe) *SearchService {
	return &SearchService{
		recipes: recipes,
	}
}

// SearchRecipes searches for recipes using string similarity matching
func (ss *SearchService) SearchRecipes(query string) []Recipe {
	if query == "" {
		return []Recipe{}
	}

	queryLower := strings.ToLower(query)
	var recipeMatches []recipeMatch
	seenNames := make(map[string]bool) // Track seen recipe names for deduplication

	// Search through recipes and calculate similarity scores
	for _, recipe := range ss.recipes {
		score := ss.calculateSimilarityScore(queryLower, recipe)
		if score > 0 {
			// Skip duplicates based on recipe name
			if seenNames[recipe.Name] {
				continue
			}
			seenNames[recipe.Name] = true

			recipeMatches = append(recipeMatches, recipeMatch{
				recipe: *recipe,
				score:  score,
			})
		}
	}

	// Sort by similarity score first, then by rating as tiebreaker
	sort.Slice(recipeMatches, func(i, j int) bool {
		// If scores are very close (within 5 points), use rating as tiebreaker
		if abs(recipeMatches[i].score-recipeMatches[j].score) < 5.0 {
			return recipeMatches[i].recipe.Rating > recipeMatches[j].recipe.Rating
		}
		return recipeMatches[i].score > recipeMatches[j].score
	})

	// Convert to Recipe slice and limit to top 20 results
	var results []Recipe
	limit := 20
	if len(recipeMatches) < limit {
		limit = len(recipeMatches)
	}

	for i := 0; i < limit; i++ {
		results = append(results, recipeMatches[i].recipe)
	}

	return results
}

// recipeMatch holds a recipe with its similarity score
type recipeMatch struct {
	recipe Recipe
	score  float64
}

// calculateSimilarityScore calculates how similar a recipe is to the search query
func (ss *SearchService) calculateSimilarityScore(query string, recipe *Recipe) float64 {
	var score float64

	// Normalize strings for comparison
	recipeName := strings.ToLower(recipe.Name)
	recipeIngredients := strings.ToLower(strings.Join(recipe.Ingredients, " "))
	recipeCuisine := strings.ToLower(recipe.CuisinePath)

	// Exact match in recipe name gets highest score
	if strings.Contains(recipeName, query) {
		score += 100.0
		// Bonus for exact word match
		if ss.containsWholeWord(recipeName, query) {
			score += 50.0
		}
	}

	// Partial matches in recipe name
	queryWords := strings.Fields(query)
	for _, word := range queryWords {
		if len(word) > 2 && strings.Contains(recipeName, word) {
			score += 30.0
			if ss.containsWholeWord(recipeName, word) {
				score += 20.0
			}
		}
	}

	// Matches in ingredients
	if strings.Contains(recipeIngredients, query) {
		score += 40.0
	}
	for _, word := range queryWords {
		if len(word) > 2 && strings.Contains(recipeIngredients, word) {
			score += 15.0
		}
	}

	// Matches in cuisine path
	if strings.Contains(recipeCuisine, query) {
		score += 25.0
	}
	for _, word := range queryWords {
		if len(word) > 2 && strings.Contains(recipeCuisine, word) {
			score += 10.0
		}
	}

	// Fuzzy matching for typos and similar words
	score += ss.calculateFuzzyScore(query, recipeName) * 20.0
	
	// Bonus for recipes with higher ratings (more significant weight)
	if recipe.Rating > 0 {
		score += recipe.Rating * 10.0 // Increased weight for ratings
	}

	return score
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// containsWholeWord checks if a word appears as a complete word (not as part of another word)
func (ss *SearchService) containsWholeWord(text, word string) bool {
	words := strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
	
	for _, w := range words {
		if w == word {
			return true
		}
	}
	return false
}

// calculateFuzzyScore calculates a fuzzy similarity score between two strings
func (ss *SearchService) calculateFuzzyScore(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	
	// Simple character-based similarity
	longer := s1
	shorter := s2
	if len(s2) > len(s1) {
		longer = s2
		shorter = s1
	}
	
	if len(longer) == 0 {
		return 0.0
	}
	
	// Count common characters
	commonChars := 0
	for _, char := range shorter {
		if strings.ContainsRune(longer, char) {
			commonChars++
		}
	}
	
	return float64(commonChars) / float64(len(longer))
}

// GetStats returns statistics about the recipe collection
func (ss *SearchService) GetStats() (int, int) {
	totalRecipes := len(ss.recipes)
	recipesWithRating := 0

	for _, recipe := range ss.recipes {
		if recipe.Rating > 0 {
			recipesWithRating++
		}
	}

	return totalRecipes, recipesWithRating
}
