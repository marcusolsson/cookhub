package main

import (
	"fmt"
	"net/http"
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
