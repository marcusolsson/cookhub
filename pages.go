package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
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

	fmt.Println(fileRef)

	ctx := req.Context()

	if item, found := s.c.Get(fileRef.ID()); found {
		return item.(templ.Component).Render(ctx, w)
	}

	file, err := s.db.GetFile(ctx, db.GetFileParams{
		Provider: repoRef.Provider,
		Owner:    repoRef.Owner,
		RepoName: repoRef.Name,
		Path:     fileRef.Path,
	})
	if err != nil {
		return err
	}

	recipe, err := parseCooklangRecipe(file.Content)
	if err != nil {
		return err
	}
	model := &views.RecipeViewModel{
		Recipe:    recipe,
		RawRecipe: file.Content,
		File:      fileRef,
	}

	page := views.RecipePage(model)

	s.c.Set(fileRef.ID(), page, cache.DefaultExpiration)

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

func (s *Server) pageListCookbooks(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	repos, err := s.db.ListRepositories(ctx)
	if err != nil {
		return fmt.Errorf("failed to get recipes from database: %w", err)
	}

	return views.AllCookbooksPage(repos).Render(ctx, w)
}
