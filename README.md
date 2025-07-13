# Recipe API

A high-performance Go-based REST API for searching and retrieving recipe data with intelligent ranking based on both ratings and review counts.

## Features

- **Title-based Search**: Searches recipes strictly by title/name for precise matching
- **Intelligent Ranking**: Uses a weighted scoring algorithm that balances rating quality with review volume
- **High-Quality Results**: Only returns recipes with ratings ≥ 4.0 and at least one review
- **Fast Performance**: Efficiently handles large datasets (200K+ recipes, 1M+ interactions)
- **RESTful API**: Clean JSON responses with comprehensive metadata

## Architecture

The application is structured with clean separation of concerns:

```
├── main.go           # Application entry point and server setup
├── models.go         # Data structures and types
├── data_loader.go    # CSV data loading and parsing
├── search_service.go # Recipe search logic and ranking
├── handlers.go       # HTTP request handlers
├── go.mod           # Go module dependencies
├── RAW_recipes.csv  # Recipe data (231K+ recipes)
└── RAW_interactions.csv # User interaction data (1M+ interactions)
```

## API Endpoints

### Search Recipes
```
GET /recipes?query={search_term}
```

**Parameters:**
- `query` (required): Search term to match against recipe titles

**Response:**
```json
{
  "recipes": [
    {
      "id": 39087,
      "name": "creamy cajun chicken pasta",
      "description": "n'awlin's style of chicken with an updated alfredo sauce.",
      "minutes": 25,
      "ingredients": ["boneless skinless chicken breast halves", "linguine", ...],
      "steps": ["place chicken and cajun seasoning in a bowl...", ...],
      "avg_rating": 4.541436464088398,
      "n_reviews": 1448
    }
  ],
  "count": 20,
  "query": "chicken"
}
```

**Example Requests:**
```bash
curl "http://localhost:8080/recipes?query=chicken"
curl "http://localhost:8080/recipes?query=pasta"
curl "http://localhost:8080/recipes?query=chocolate"
```

### Health Check
```
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "total_recipes": 231637,
  "good_recipes": 194707,
  "total_interactions": 1132367
}
```

## Ranking Algorithm

The API uses a sophisticated weighted scoring system that considers both rating quality and review volume:

```
Score = Average Rating × log₁₀(Number of Reviews + 1)
```

**Why this works:**
- A recipe with 2000 reviews at 4.8★ scores ~15.8
- A recipe with 3 reviews at 5.0★ scores ~3.0
- Prevents recipes with few reviews from dominating results
- Gives appropriate weight to well-reviewed recipes

## Data Processing

### Recipe Data
- **Source**: `RAW_recipes.csv` (231,637 recipes)
- **Fields**: ID, name, description, cooking time, ingredients, steps
- **Processing**: Parses list-formatted strings, handles malformed records

### Interaction Data
- **Source**: `RAW_interactions.csv` (1,132,367 interactions)
- **Fields**: User ID, recipe ID, date, rating (1-5), review text
- **Processing**: Calculates average ratings and review counts per recipe

### Quality Filtering
- Only recipes with ≥1 review are considered
- Only recipes with ≥4.0 average rating are returned
- Results in ~195K high-quality recipes available for search

## Getting Started

### Prerequisites
- Go 1.19 or higher
- CSV data files (`RAW_recipes.csv`, `RAW_interactions.csv`)

### Installation
```bash
# Clone the repository
git clone <repository-url>
cd recipe-api

# Run the application
go run .
```

### Usage
```bash
# Start the server
go run .

# Search for chicken recipes
curl "http://localhost:8080/recipes?query=chicken"

# Check API health
curl "http://localhost:8080/health"
```

## Performance

- **Startup Time**: ~2-3 seconds (loads 231K recipes + 1M interactions)
- **Search Response**: <100ms for typical queries
- **Memory Usage**: ~200MB for full dataset
- **Concurrent Requests**: Supports high concurrency with Go's built-in HTTP server

## Development

### Project Structure
- **models.go**: Defines `Recipe`, `Interaction`, and `RecipeWithStats` structs
- **data_loader.go**: Handles CSV parsing and data loading with progress logging
- **search_service.go**: Implements search logic and weighted ranking algorithm
- **handlers.go**: HTTP request handlers with JSON response formatting
- **main.go**: Application bootstrap and server configuration

### Key Design Decisions
1. **In-Memory Storage**: All data loaded into memory for fast search performance
2. **Title-Only Search**: Focused search scope for more relevant results
3. **Weighted Scoring**: Balances rating quality with review confidence
4. **Modular Architecture**: Clean separation for maintainability and testing

### Adding Features
- **New Endpoints**: Add handlers in `handlers.go`
- **Search Improvements**: Modify `search_service.go`
- **Data Sources**: Extend `data_loader.go`
- **Response Format**: Update models in `models.go`

## Example Queries

### Popular Chicken Recipes
```bash
curl "http://localhost:8080/recipes?query=chicken" | jq '.recipes[0:3]'
```

### Pasta Dishes
```bash
curl "http://localhost:8080/recipes?query=pasta" | jq '.recipes[0:5]'
```

### Desserts
```bash
curl "http://localhost:8080/recipes?query=cake" | jq '.recipes[0:3]'
```

## Monitoring

The `/health` endpoint provides key metrics:
- **total_recipes**: All loaded recipes
- **good_recipes**: High-quality recipes (≥4.0 rating, ≥1 review)
- **total_interactions**: All user interactions processed
- **status**: API health status

## Future Enhancements

- [ ] Add recipe filtering by cooking time, ingredients, etc.
- [ ] Implement recipe recommendation based on user preferences
- [ ] Add caching layer for frequently searched terms
- [ ] Support for fuzzy/typo-tolerant search
- [ ] Add recipe categories and tags
- [ ] Implement rate limiting and API authentication
- [ ] Add comprehensive logging and metrics
- [ ] Support for recipe image URLs
