package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// DataLoader handles loading and processing of CSV data
type DataLoader struct {
	recipes      map[int]*RecipeWithStats
	interactions []Interaction
}

// NewDataLoader creates a new DataLoader instance
func NewDataLoader() *DataLoader {
	return &DataLoader{
		recipes: make(map[int]*RecipeWithStats),
	}
}

// LoadData loads both recipes and interactions from CSV files
func (dl *DataLoader) LoadData() error {
	log.Println("Loading recipes from CSV...")
	if err := dl.loadRecipes(); err != nil {
		return fmt.Errorf("failed to load recipes: %v", err)
	}

	log.Println("Loading interactions from CSV...")
	if err := dl.loadInteractions(); err != nil {
		return fmt.Errorf("failed to load interactions: %v", err)
	}

	log.Println("Data loading completed successfully")
	return nil
}

// loadRecipes loads recipe data from RAW_recipes.csv
func (dl *DataLoader) loadRecipes() error {
	file, err := os.Open("RAW_recipes.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header row
	for i, record := range records[1:] {
		if len(record) < 12 {
			log.Printf("Skipping malformed recipe record at line %d", i+2)
			continue
		}

		id, err := strconv.Atoi(record[1])
		if err != nil {
			log.Printf("Invalid recipe ID at line %d: %v", i+2, err)
			continue
		}

		minutes, err := strconv.Atoi(record[2])
		if err != nil {
			minutes = 0 // Default to 0 if parsing fails
		}

		recipe := &RecipeWithStats{
			Recipe: Recipe{
				ID:          id,
				Name:        record[0],
				Description: record[9],
				Minutes:     minutes,
				Ingredients: parseListString(record[10]),
				Steps:       parseListString(record[8]),
			},
		}

		dl.recipes[id] = recipe

		if (i+1)%1000 == 0 {
			log.Printf("Loaded %d recipes...", i+1)
		}
	}

	log.Printf("Total recipes loaded: %d", len(dl.recipes))
	return nil
}

// loadInteractions loads interaction data from RAW_interactions.csv
func (dl *DataLoader) loadInteractions() error {
	file, err := os.Open("RAW_interactions.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header row
	for i, record := range records[1:] {
		if len(record) < 5 {
			log.Printf("Skipping malformed interaction record at line %d", i+2)
			continue
		}

		userID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("Invalid user ID at line %d: %v", i+2, err)
			continue
		}

		recipeID, err := strconv.Atoi(record[1])
		if err != nil {
			log.Printf("Invalid recipe ID at line %d: %v", i+2, err)
			continue
		}

		rating, err := strconv.Atoi(record[3])
		if err != nil {
			log.Printf("Invalid rating at line %d: %v", i+2, err)
			continue
		}

		interaction := Interaction{
			UserID:   userID,
			RecipeID: recipeID,
			Date:     record[2],
			Rating:   rating,
			Review:   record[4],
		}

		dl.interactions = append(dl.interactions, interaction)

		if (i+1)%5000 == 0 {
			log.Printf("Loaded %d interactions...", i+1)
		}
	}

	log.Printf("Total interactions loaded: %d", len(dl.interactions))
	return nil
}

// CalculateRecipeStats calculates average ratings and review counts for all recipes
func (dl *DataLoader) CalculateRecipeStats() {
	log.Println("Calculating recipe statistics...")

	for _, interaction := range dl.interactions {
		if recipe, exists := dl.recipes[interaction.RecipeID]; exists {
			recipe.TotalRating += interaction.Rating
			recipe.ReviewCount++
		}
	}

	// Calculate average ratings
	for _, recipe := range dl.recipes {
		if recipe.ReviewCount > 0 {
			recipe.AvgRating = float64(recipe.TotalRating) / float64(recipe.ReviewCount)
			recipe.NReviews = recipe.ReviewCount
		}
	}

	log.Println("Recipe statistics calculated successfully")
}

// GetRecipes returns the loaded recipes map
func (dl *DataLoader) GetRecipes() map[int]*RecipeWithStats {
	return dl.recipes
}

// GetInteractions returns the loaded interactions slice
func (dl *DataLoader) GetInteractions() []Interaction {
	return dl.interactions
}

// parseListString parses a string representation of a list into a slice of strings
func parseListString(listStr string) []string {
	// Remove brackets and quotes, then split by comma
	listStr = strings.Trim(listStr, "[]")
	if listStr == "" {
		return []string{}
	}

	// Split by comma and clean each item
	items := strings.Split(listStr, "', '")
	var result []string
	for _, item := range items {
		cleaned := strings.Trim(item, "'\"")
		if cleaned != "" {
			result = append(result, cleaned)
		}
	}
	return result
}
