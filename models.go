package main

// Recipe represents a recipe with all its details from the new CSV format
type Recipe struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	PrepTime    string   `json:"prep_time"`
	CookTime    string   `json:"cook_time"`
	TotalTime   string   `json:"total_time"`
	Servings    string   `json:"servings"`
	Yield       string   `json:"yield"`
	Ingredients []string `json:"ingredients"`
	Directions  []string `json:"directions"`
	Rating      float64  `json:"rating"`
	URL         string   `json:"url"`
	CuisinePath string   `json:"cuisine_path"`
	Nutrition   string   `json:"nutrition"`
	Timing      string   `json:"timing"`
	ImageSrc    string   `json:"img_src"`
}
