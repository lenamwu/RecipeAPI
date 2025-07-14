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
	recipes map[int]*Recipe
}

// NewDataLoader creates a new DataLoader instance
func NewDataLoader() *DataLoader {
	return &DataLoader{
		recipes: make(map[int]*Recipe),
	}
}

// LoadData loads recipes from the new CSV file
func (dl *DataLoader) LoadData() error {
	log.Println("Loading recipes from recipes.csv...")
	if err := dl.loadRecipes(); err != nil {
		return fmt.Errorf("failed to load recipes: %v", err)
	}

	log.Println("Data loading completed successfully")
	return nil
}

// loadRecipes loads recipe data from recipes.csv
func (dl *DataLoader) loadRecipes() error {
	file, err := os.Open("recipes.csv")
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
		if len(record) < 15 {
			log.Printf("Skipping malformed recipe record at line %d", i+2)
			continue
		}

		// Parse ID from the first column (index)
		id, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("Invalid recipe ID at line %d: %v", i+2, err)
			continue
		}

		// Parse rating
		rating := 0.0
		if record[9] != "" {
			if r, err := strconv.ParseFloat(record[9], 64); err == nil {
				rating = r
			}
		}

		// Parse ingredients and directions
		ingredients := parseIngredients(record[7])
		directions := parseDirections(record[8])

		recipe := &Recipe{
			ID:          id,
			Name:        record[1],
			PrepTime:    record[2],
			CookTime:    record[3],
			TotalTime:   record[4],
			Servings:    record[5],
			Yield:       record[6],
			Ingredients: ingredients,
			Directions:  directions,
			Rating:      rating,
			URL:         record[10],
			CuisinePath: record[11],
			Nutrition:   record[12],
			Timing:      record[13],
			ImageSrc:    record[14],
		}

		dl.recipes[id] = recipe

		if (i+1)%1000 == 0 {
			log.Printf("Loaded %d recipes...", i+1)
		}
	}

	log.Printf("Total recipes loaded: %d", len(dl.recipes))
	return nil
}

// GetRecipes returns the loaded recipes map
func (dl *DataLoader) GetRecipes() map[int]*Recipe {
	return dl.recipes
}

// parseIngredients parses the ingredients string into a slice
func parseIngredients(ingredientsStr string) []string {
	if ingredientsStr == "" {
		return []string{}
	}
	
	// Split by comma and clean each ingredient
	ingredients := strings.Split(ingredientsStr, ",")
	var result []string
	for _, ingredient := range ingredients {
		cleaned := strings.TrimSpace(ingredient)
		if cleaned != "" {
			result = append(result, cleaned)
		}
	}
	return result
}

// parseDirections parses the directions string into a slice
func parseDirections(directionsStr string) []string {
	if directionsStr == "" {
		return []string{}
	}
	
	// Split by sentence endings and clean each direction
	directions := strings.Split(directionsStr, ". ")
	var result []string
	for _, direction := range directions {
		cleaned := strings.TrimSpace(direction)
		if cleaned != "" {
			// Add period back if it was removed by split
			if !strings.HasSuffix(cleaned, ".") && !strings.HasSuffix(cleaned, "!") && !strings.HasSuffix(cleaned, "?") {
				cleaned += "."
			}
			result = append(result, cleaned)
		}
	}
	return result
}
