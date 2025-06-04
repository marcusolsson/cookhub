package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
)

type server struct {
	db *db.Queries
}

func newServer(qs *db.Queries) chi.Router {
	srv := &server{
		db: qs,
	}
	r := chi.NewRouter()
	r.Get("/github.com/{org}/{repo}/*", srv.getRecipe)
	r.Get("/jobs", srv.getJobs)
	r.Get("/recipes", srv.getRecipes)

	r.Route("/api", func(r chi.Router) {
		r.Post("/github.com/{org}/{name}", srv.indexRepo)
		r.Get("/github.com/{org}/{name}", srv.getIngestedFiles)
	})
	return r
}

func (s *server) getRecipe(w http.ResponseWriter, req *http.Request) {
	var (
		org  = chi.URLParam(req, "org")
		repo = chi.URLParam(req, "repo")
		slug = fmt.Sprintf("%s/%s", org, repo)
	)

	relPath := strings.TrimPrefix(
		req.URL.Path,
		fmt.Sprintf("/ui/github.com/%s/%s", org, repo),
	)

	ctx := req.Context()

	jobID, err := s.db.GetLatestJobBySlug(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	file, err := s.db.GetFileByName(ctx, db.GetFileByNameParams{
		JobID: jobID,
		Name:  relPath,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	cooklangRecipe, err := ParseCooklangRecipe(file.Name, file.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	schemaOrgRecipe := ConvertCooklangToSchemaOrg(file.Name, cooklangRecipe)

	component := recipeView(schemaOrgRecipe)
	component.Render(req.Context(), w)
}

func (s *server) getJobs(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	jobs, err := s.db.GetJobs(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	component := jobsView(jobs)
	component.Render(req.Context(), w)
}

func (s *server) getRecipes(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	recipes, err := s.db.GetRecipes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	component := recipesView(recipes)
	component.Render(req.Context(), w)
}

func (s *server) getIngestedFiles(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var (
		org  = chi.URLParam(req, "org")
		name = chi.URLParam(req, "name")
		slug = fmt.Sprintf("%s/%s", org, name)
	)

	jobID, err := s.db.GetLatestJobBySlug(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files, err := s.db.GetFilesByJob(ctx, jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) indexRepo(w http.ResponseWriter, req *http.Request) {
	var (
		org      = chi.URLParam(req, "org")
		repoName = chi.URLParam(req, "name")
		slug     = fmt.Sprintf("%s/%s", org, repoName)
	)

	ctx := req.Context()

	zipPath, err := downloadZipBall(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the commit SHA from the zip filepath.
	sha := strings.Split(
		strings.TrimSuffix(
			filepath.Base(zipPath),
			filepath.Ext(zipPath),
		),
		"-",
	)[2]

	jobID, err := s.db.CreateJob(ctx, db.CreateJobParams{Slug: slug, CommitSha: sha})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for file := range readFilesFromZip(zipPath) {
		if err := s.db.CreateFile(ctx, db.CreateFileParams{
			JobID:   jobID,
			Name:    file.Name,
			Content: string(file.Content),
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
