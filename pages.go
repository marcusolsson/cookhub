package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aquilax/cooklang-go"
	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
	"github.com/marcusolsson/cookhub/github"
	"github.com/marcusolsson/cookhub/utils"
	"github.com/marcusolsson/cookhub/views"
	"github.com/patrickmn/go-cache"
)

type recipeViewModel struct {
	metadata    *utils.RecipeMetadata
	ghRepo      github.Repository
	recipe      *cooklang.RecipeV2
	fileContent string
}

// pageShowRecipe renders the page for a specific recipe.
func (s *Server) pageShowRecipe(w http.ResponseWriter, req *http.Request) error {
	repoRef := utils.RepoRef{
		Provider: chi.URLParam(req, "provider"),
		Owner:    chi.URLParam(req, "owner"),
		Name:     chi.URLParam(req, "name"),
	}

	fileRef := utils.RepoFileRef{
		Repo: repoRef,
		Path: chi.URLParam(req, "*"),
		Ref:  "HEAD",
	}

	ctx := req.Context()

	ctxlog := s.logger.With("file", fileRef.ID())

	if item, found := s.c.Get(fileRef.ID()); found {
		vm := item.(recipeViewModel)

		component := views.RecipeView(
			vm.metadata,
			vm.ghRepo,
			vm.recipe,
			vm.fileContent,
			fileRef,
		)

		return component.Render(ctx, w)
	}

	jobID, err := s.db.GetLatestJobBySlug(ctx, repoRef.Slug())
	if err != nil {
		return fmt.Errorf("failed to get latest job by slug: %w", err)
	}

	ctxlog = ctxlog.With("job_id", jobID)

	repoMetadata, err := s.db.GetRepoMetadataByJob(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get repo metadata: %w", err)
	}

	var repo github.Repository
	json.Unmarshal(repoMetadata.Response, &repo)

	file, err := s.db.GetFileByName(ctx, db.GetFileByNameParams{
		JobID: jobID,
		Name:  fileRef.Path,
	})
	if err != nil {
		return fmt.Errorf("failed to get file by name: %w", err)
	}

	recipe, err := parseCooklangRecipe(file.Content)
	if err != nil {
		return fmt.Errorf("failed to parse Cooklang recipe: %w", err)
	}

	metadata := utils.ParseCanonicalMetadata(recipe, file.Name)

	s.c.Set(fileRef.ID(), recipeViewModel{
		metadata:    metadata,
		ghRepo:      repo,
		recipe:      recipe,
		fileContent: file.Content,
	}, cache.DefaultExpiration)

	component := views.RecipeView(metadata, repo, recipe, file.Content, fileRef)

	return component.Render(ctx, w)
}

// pageListJobs renders the page that lists all jobs.
func (s *Server) pageListJobs(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	jobs, err := s.db.GetJobs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get jobs from database: %w", err)
	}

	return views.JobsView(jobs).Render(ctx, w)
}

// pageListRecipes renders the page that lists all recipes.
func (s *Server) pageListRecipes(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	recipes, err := s.db.GetRecipes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get recipes from database: %w", err)
	}

	return views.RecipesView(recipes).Render(ctx, w)
}
