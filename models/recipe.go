package models

// RecipeFile represents a single .cook file for static generation
type RecipeFile struct {
	// Filesystem-derived fields
	Path      string // relative path within input directory (e.g., "breakfast/pancakes.cook")
	Stem      string // filename without extension (e.g., "pancakes")
	Content   string // raw file content

	// Output-related fields
	OutputPath string // relative path for HTML output (e.g., "breakfast/pancakes.html")
	URL        string // URL path for linking (e.g., "/breakfast/pancakes.html")
}

// Cookbook represents the collection of recipes being generated
type Cookbook struct {
	Title   string       // site title
	BaseURL string       // base URL path prefix (e.g., "/recipes" for serving at example.com/recipes/)
	Recipes []RecipeFile // all discovered recipes
}
