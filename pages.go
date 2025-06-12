package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
	"github.com/marcusolsson/cookhub/utils"
	"github.com/marcusolsson/cookhub/views"
)

// pageShowRecipe renders the page for a specific recipe.
func (s *Server) pageShowRecipe(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	file, err := s.db.GetFile(ctx, db.GetFileParams{
		Provider: chi.URLParam(req, "provider"),
		Owner:    chi.URLParam(req, "owner"),
		RepoName: chi.URLParam(req, "name"),
		Path:     chi.URLParam(req, "*"),
	})
	if err != nil {
		return err
	}

	recipe, err := parseCooklangRecipe(file.Content)
	if err != nil {
		return err
	}
	model := &views.RecipeViewModel{
		Recipe: recipe,
		File:   file,
	}

	page := views.RecipePage(model)

	return page.Render(ctx, w)
}

func (s *Server) pageListRecipesByRepo(w http.ResponseWriter, req *http.Request) error {
	repoRef := utils.RepoRef{
		Provider: chi.URLParam(req, "provider"),
		Owner:    chi.URLParam(req, "owner"),
		Name:     chi.URLParam(req, "name"),
	}

	ctx := req.Context()

	files, err := s.db.GetFiles(ctx, db.GetFilesParams{
		Provider: repoRef.Provider,
		Owner:    repoRef.Owner,
		RepoName: repoRef.Name,
	})
	if err != nil {
		return fmt.Errorf("failed to get recipes from database: %w", err)
	}

	return views.AllRecipesPage(repoRef, files).Render(ctx, w)
}
