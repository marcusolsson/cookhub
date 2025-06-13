package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
	"github.com/marcusolsson/cookhub/views"
)

// recipePage renders the page for a specific recipe.
func (s *Server) recipePage(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	params := db.GetFileParams{
		Provider: chi.URLParam(req, "provider"),
		Owner:    chi.URLParam(req, "owner"),
		RepoName: chi.URLParam(req, "repo_name"),
		Path:     chi.URLParam(req, "*"),
	}

	file, err := s.db.GetFile(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			msg := fmt.Sprintf(
				"We couldn't find a recipe called %q in %s/%s.",
				params.Path,
				params.Owner,
				params.RepoName,
			)

			return views.NotFoundPage(msg).Render(ctx, w)
		}
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

func (s *Server) cookbookPage(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	params := db.FindRepositoryByNameParams{
		Provider: chi.URLParam(req, "provider"),
		Owner:    chi.URLParam(req, "owner"),
		RepoName: chi.URLParam(req, "repo_name"),
	}

	repo, err := s.db.FindRepositoryByName(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			msg := fmt.Sprintf(
				"We couldn't find a cookbook called %q by %s.",
				params.RepoName,
				params.Owner,
			)

			return views.NotFoundPage(msg).Render(ctx, w)
		}
		return fmt.Errorf("failed to get recipes from database: %w", err)
	}

	files, err := s.db.GetFiles(ctx, db.GetFilesParams{
		Provider: repo.Provider,
		Owner:    repo.Owner,
		RepoName: repo.RepoName,
	})
	if err != nil {
		return fmt.Errorf("failed to get recipes from database: %w", err)
	}

	return views.AllRecipesPage(repo, files).Render(ctx, w)
}

func (s *Server) notFoundPage(w http.ResponseWriter, req *http.Request) error {
	return views.NotFoundPage("We couldn't find a page like that.").Render(req.Context(), w)
}
