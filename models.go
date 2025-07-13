package main

// Recipe represents a recipe with all its details
type Recipe struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Minutes     int      `json:"minutes"`
	Ingredients []string `json:"ingredients"`
	Steps       []string `json:"steps"`
	AvgRating   float64  `json:"avg_rating"`
	NReviews    int      `json:"n_reviews"`
}

// Interaction represents a user interaction with a recipe
type Interaction struct {
	UserID   int    `json:"user_id"`
	RecipeID int    `json:"recipe_id"`
	Date     string `json:"date"`
	Rating   int    `json:"rating"`
	Review   string `json:"review"`
}

// RecipeWithStats extends Recipe with additional statistics for internal processing
type RecipeWithStats struct {
	Recipe
	TotalRating int
	ReviewCount int
}
