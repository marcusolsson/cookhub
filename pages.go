package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
)

func (s *server) pageShowRecipe(w http.ResponseWriter, req *http.Request) {
	var (
		org  = chi.URLParam(req, "org")
		repo = chi.URLParam(req, "repo")
		slug = fmt.Sprintf("%s/%s", org, repo)
	)

	relPath := strings.TrimPrefix(
		req.URL.Path,
		fmt.Sprintf("/github.com/%s/%s", org, repo),
	)

	ctx := req.Context()

	ctxlog := s.logger.With("slug", slug, "file_path", relPath)

	jobID, err := s.db.GetLatestJobBySlug(ctx, slug)
	if err != nil {
		ctxlog.Error("Failed to get latest job by slug", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctxlog = ctxlog.With("job_id", jobID)

	file, err := s.db.GetFileByName(ctx, db.GetFileByNameParams{
		JobID: jobID,
		Name:  relPath,
	})
	if err != nil {
		ctxlog.Error("Failed to get file by name", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cooklangRecipe, err := ParseCooklangRecipe(file.Name, file.Content)
	if err != nil {
		ctxlog.Error("Failed to parse Cooklang recipe", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	schemaOrgRecipe := ConvertCooklangToSchemaOrg(file.Name, cooklangRecipe)

	recipeView(schemaOrgRecipe).Render(ctx, w)
}

func (s *server) pageListJobs(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	jobs, err := s.db.GetJobs(ctx)
	if err != nil {
		s.logger.Error("Failed to get jobs from database", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jobsView(jobs).Render(ctx, w)
}

func (s *server) pageListRecipes(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	recipes, err := s.db.GetRecipes(ctx)
	if err != nil {
		s.logger.Error("Failed to get recipes from database", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recipesView(recipes).Render(ctx, w)
}
