package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aquilax/cooklang-go"
	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
	"github.com/marcusolsson/cookhub/views"
	"github.com/patrickmn/go-cache"
)

type recipeViewModel struct {
	metadata    *views.RecipeMetadata
	author      views.Author
	recipe      *cooklang.RecipeV2
	fileContent string
}

// pageShowRecipe renders the page for a specific recipe.
func (s *server) pageShowRecipe(w http.ResponseWriter, req *http.Request) error {
	var (
		provider = chi.URLParam(req, "provider")
		owner    = chi.URLParam(req, "owner")
		name     = chi.URLParam(req, "name")
		filename = chi.URLParam(req, "*")
		fullName = fmt.Sprintf("%s/%s", owner, name)
	)

	fullURL := strings.Join([]string{
		provider,
		owner,
		name,
		filename,
	},
		"/",
	)

	ctx := req.Context()

	ctxlog := s.logger.With("slug", fullName, "filename", filename)

	if item, found := s.c.Get(fullURL); found {
		vm := item.(recipeViewModel)

		component := views.RecipeView(
			vm.metadata,
			vm.author,
			vm.recipe,
			vm.fileContent,
		)

		return component.Render(ctx, w)
	}

	filename = "/" + filename

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

	s.c.Set(fullURL, recipeViewModel{
		metadata:    metadata,
		author:      author,
		recipe:      recipe,
		fileContent: file.Content,
	}, cache.DefaultExpiration)

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
