package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
	"github.com/marcusolsson/cookhub/views"
)

// pageShowRecipe renders the page for a specific recipe.
func (s *server) pageShowRecipe(w http.ResponseWriter, req *http.Request) error {
	var (
		org  = chi.URLParam(req, "org")
		repo = chi.URLParam(req, "repo")
		slug = fmt.Sprintf("%s/%s", org, repo)
	)

	filename := strings.TrimPrefix(
		req.URL.Path,
		fmt.Sprintf("/github.com/%s/%s", org, repo),
	)

	ctx := req.Context()

	ctxlog := s.logger.With("slug", slug, "filename", filename)

	jobID, err := s.db.GetLatestJobBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("failed to get latest job by slug: %w", err)
	}

	ctxlog = ctxlog.With("job_id", jobID)

	file, err := s.db.GetFileByName(ctx, db.GetFileByNameParams{
		JobID: jobID,
		Name:  filename,
	})
	if err != nil {
		return fmt.Errorf("failed to get file by name: %w", err)
	}

	recipe, err := parseCooklangRecipe(file.Content)
	if err != nil {
		return fmt.Errorf("failed to parse Cooklang recipe: %w", err)
	}

	name := strings.TrimSuffix(
		filepath.Base(file.Name),
		filepath.Ext(file.Name),
	)

	return views.RecipeView(name, recipe).Render(ctx, w)
}

// pageListJobs renders the page that lists all jobs.
func (s *server) pageListJobs(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	jobs, err := s.db.GetJobs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get jobs from database: %w", err)
	}

	return views.JobsView(jobs).Render(ctx, w)
}

// pageListRecipes renders the page that lists all recipes.
func (s *server) pageListRecipes(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	recipes, err := s.db.GetRecipes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get recipes from database: %w", err)
	}

	return views.RecipesView(recipes).Render(ctx, w)
}
