package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
)

type server struct {
	db *db.Queries

	githubClient *GitHubClient
}

func newServer(qs *db.Queries) chi.Router {
	srv := &server{
		db:           qs,
		githubClient: newGitHubClient(),
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

	sha, err := s.githubClient.GetLatestCommitSHA(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jobID, err := s.db.CreateJob(ctx, db.CreateJobParams{Slug: slug, CommitSha: sha})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = func() error {
		zipPath, err := downloadZipBall(ctx, slug, sha)
		if err != nil {
			return err
		}

		var params []db.CreateFileParams
		for file := range readFilesFromZip(zipPath) {
			params = append(params, db.CreateFileParams{
				JobID:   jobID,
				Name:    file.Name,
				Content: string(file.Content),
			})
		}

		if _, err := s.db.CreateFile(ctx, params); err != nil {
			return err
		}

		return nil
	}()
	if err != nil {
		s.db.SetJobStatus(ctx, db.SetJobStatusParams{
			Status: db.StatusEnumFailed,
			ID:     jobID,
		})
		return
	}

	s.db.SetJobStatus(ctx, db.SetJobStatusParams{
		Status: db.StatusEnumCompleted,
		ID:     jobID,
	})
}
