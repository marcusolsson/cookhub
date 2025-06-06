package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
	"github.com/marcusolsson/cookhub/views"
)

// pageShowRecipe renders the page for a specific recipe.
func (s *server) pageShowRecipe(w http.ResponseWriter, req *http.Request) error {
	var (
		owner    = chi.URLParam(req, "owner")
		name     = chi.URLParam(req, "name")
		filename = "/" + chi.URLParam(req, "*")
		fullName = fmt.Sprintf("%s/%s", owner, name)
	)

	ctx := req.Context()

	ctxlog := s.logger.With("slug", fullName, "filename", filename)

	jobID, err := s.db.GetLatestJobBySlug(ctx, fullName)
	if err != nil {
		return fmt.Errorf("failed to get latest job by slug: %w", err)
	}

	ctxlog = ctxlog.With("job_id", jobID)

	repo, err := s.db.GetRepoMetadataByJob(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get repo metadata: %w", err)
	}

	var author views.Author
	json.Unmarshal(repo.Response, &author)

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

	metadata := views.ParseCanonicalMetadata(recipe, file)
	component := views.RecipeView(metadata, author, recipe, file.Content)

	return component.Render(ctx, w)
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
