# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

- `go build ./cmd/cookhub` - Build the CLI
- `go test ./...` - Run all tests
- `go run github.com/a-h/templ/cmd/templ@latest generate` - Generate Go code from Templ templates

## Usage

```bash
cookhub generate --input <recipes-dir> --output <output-dir> [--title "Site Title"] [--base-url "/prefix"]
```

Flags:

- `--input` (required): Directory containing `.cook` recipe files
- `--output`: Output directory for generated HTML (default: "public")
- `--title`: Site title for the index page
- `--base-url`: Base URL path prefix for serving site from a subpath

## Architecture

CookHub is a static site generator CLI that converts CookLang (.cook) recipe files to HTML.

**Tech Stack:** Go 1.23.2, Templ (templates), CookLang parser, Cobra (CLI)

**Generation Flow:**
1. CLI (`cmd/cookhub/`) parses flags and invokes generator
2. Generator (`generator/generator.go`) discovers `.cook` files recursively
3. CookLang recipes parsed with `cooklang-go` library
4. HTML rendered with Templ components in `views/`
5. Static assets copied from embedded `generator/static/`

**Packages (3 total):**

- `cmd/cookhub/` - CLI entry point (main.go, root.go, generate.go)
- `generator/` - Core generation logic and embedded static assets
- `views/` - Templ templates, data models, and metadata parsing

**Data Model (in views/):**

- `RecipeFile` - Represents a single `.cook` file with path, content, and output URL
- `Cookbook` - Collection of recipes with site title and base URL
- `RecipeMetadata` - Parsed recipe metadata (title, description, times, tags, etc.)

**Output Structure:**

```text
output/
  index.html              # Recipe index page
  recipe-name/index.html  # Individual recipe (clean URLs)
  category/recipe/index.html
  static/styles.css       # Stylesheet
  static/robots.txt       # SEO file
```
