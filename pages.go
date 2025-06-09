package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
	"github.com/marcusolsson/cookhub/github"
	"github.com/marcusolsson/cookhub/utils"
	"github.com/marcusolsson/cookhub/views"
	"github.com/patrickmn/go-cache"
)

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
		return item.(templ.Component).Render(ctx, w)
	}

	jobID, err := s.db.GetLatestJobBySlug(ctx, repoRef.Slug())
	if err != nil {
		return fmt.Errorf("failed to get latest job by slug: %w", err)
	}

	ctxlog = ctxlog.With("job_id", jobID)

	model, err := s.makeRecipeViewModel(ctx, jobID, fileRef)
	if err != nil {
		return err
	}

	page := views.RecipePage(model)

	s.c.Set(fileRef.ID(), page, cache.DefaultExpiration)

	return page.Render(ctx, w)
}

// makeRecipeViewModel constructs the view model for a recipe page.
func (s *Server) makeRecipeViewModel(
	ctx context.Context,
	jobID string,
	fileRef utils.RepoFileRef,
) (*views.RecipeViewModel, error) {
	data, err := s.db.GetRecipePageData(ctx, db.GetRecipePageDataParams{
		JobID: jobID,
		Name:  fileRef.Path,
	})
	if err != nil {
		return nil, err
	}

	var repo github.Repository
	json.Unmarshal(data.Response, &repo)

	recipe, err := parseCooklangRecipe(data.Content)
	if err != nil {
		return nil, err
	}

	return &views.RecipeViewModel{
		Repo:      repo,
		Recipe:    recipe,
		RawRecipe: data.Content,
		File:      fileRef,
	}, nil
}

// pageListRecipes renders the page that lists all recipes.
func (s *Server) pageListRecipes(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	recipes, err := s.db.GetRecipes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get recipes from database: %w", err)
	}

	return views.AllRecipesPage(recipes).Render(ctx, w)
}
