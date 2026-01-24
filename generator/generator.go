package generator

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/aquilax/cooklang-go"
	"github.com/marcusolsson/cookhub/views"
)

//go:embed static/*
var staticFiles embed.FS

// Config holds the configuration for site generation
type Config struct {
	InputDir  string
	OutputDir string
	Title     string
	BaseURL   string // base URL path prefix (e.g., "/recipes")
}

// Generate creates a static site from CookLang recipe files
func Generate(cfg Config) error {
	// Normalize base URL: ensure it starts with / and doesn't end with /
	baseURL := strings.TrimSuffix(cfg.BaseURL, "/")
	if baseURL != "" && !strings.HasPrefix(baseURL, "/") {
		baseURL = "/" + baseURL
	}

	// Discover all .cook files
	recipes, err := discoverRecipes(cfg.InputDir, baseURL)
	if err != nil {
		return fmt.Errorf("discovering recipes: %w", err)
	}

	if len(recipes) == 0 {
		return fmt.Errorf("no .cook files found in %s", cfg.InputDir)
	}

	fmt.Printf("Found %d recipe(s)\n", len(recipes))

	// Create output directory
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Generate individual recipe pages
	for _, recipe := range recipes {
		if err := generateRecipePage(cfg, recipe, baseURL); err != nil {
			return fmt.Errorf("generating recipe %s: %w", recipe.Path, err)
		}
		fmt.Printf("  Generated: %s\n", recipe.OutputPath)
	}

	// Generate index page
	cookbook := views.Cookbook{
		Title:   cfg.Title,
		BaseURL: baseURL,
		Recipes: recipes,
	}
	if err := generateIndexPage(cfg, cookbook); err != nil {
		return fmt.Errorf("generating index page: %w", err)
	}
	fmt.Println("  Generated: index.html")

	// Copy static assets
	if err := copyStaticAssets(cfg.OutputDir); err != nil {
		return fmt.Errorf("copying static assets: %w", err)
	}
	fmt.Println("  Copied: static/")

	return nil
}

// discoverRecipes walks the input directory and finds all .cook files
func discoverRecipes(inputDir string, baseURL string) ([]views.RecipeFile, error) {
	var recipes []views.RecipeFile

	err := filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".cook" {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		// Calculate relative path from input directory
		relPath, err := filepath.Rel(inputDir, path)
		if err != nil {
			return fmt.Errorf("calculating relative path: %w", err)
		}

		// Build recipe file model
		stem := strings.TrimSuffix(filepath.Base(path), ".cook")
		slug := views.Slugify(stem)
		dir := filepath.Dir(relPath)
		if dir == "." {
			dir = ""
		} else {
			dir = dir + "/"
		}
		outputPath := dir + slug + "/index.html"
		url := baseURL + "/" + dir + slug + "/"

		recipes = append(recipes, views.RecipeFile{
			Path:       relPath,
			Stem:       stem,
			Content:    string(content),
			OutputPath: outputPath,
			URL:        url,
		})

		return nil
	})

	return recipes, err
}

// parseCooklangRecipe parses a CookLang recipe from content
func parseCooklangRecipe(content string) (*cooklang.RecipeV2, error) {
	parser := cooklang.NewParserV2(&cooklang.ParseV2Config{})
	return parser.ParseString(content)
}

// generateRecipePage generates the HTML for a single recipe
func generateRecipePage(cfg Config, recipe views.RecipeFile, baseURL string) error {
	// Parse the recipe
	parsed, err := parseCooklangRecipe(recipe.Content)
	if err != nil {
		return fmt.Errorf("parsing recipe: %w", err)
	}

	// Create view model
	vm := &views.RecipeViewModel{
		Recipe:  parsed,
		File:    recipe,
		BaseURL: baseURL,
	}

	// Render to buffer
	var buf bytes.Buffer
	if err := views.RecipePage(vm).Render(context.Background(), &buf); err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	// Ensure output directory exists
	outputPath := filepath.Join(cfg.OutputDir, recipe.OutputPath)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// generateIndexPage generates the index.html listing all recipes
func generateIndexPage(cfg Config, cookbook views.Cookbook) error {
	var buf bytes.Buffer
	if err := views.AllRecipesPage(cookbook).Render(context.Background(), &buf); err != nil {
		return fmt.Errorf("rendering index template: %w", err)
	}

	outputPath := filepath.Join(cfg.OutputDir, "index.html")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing index file: %w", err)
	}

	return nil
}

// copyStaticAssets copies the embedded static files to the output directory
func copyStaticAssets(outputDir string) error {
	destDir := filepath.Join(outputDir, "static")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("creating static output directory: %w", err)
	}

	return fs.WalkDir(staticFiles, "static", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Read from embedded FS
		content, err := staticFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading embedded file %s: %w", path, err)
		}

		// Get relative path (remove "static/" prefix)
		relPath, err := filepath.Rel("static", path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)

		// Ensure destination directory exists
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(destPath, content, 0644)
	})
}
