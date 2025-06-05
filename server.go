package main

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
)

type server struct {
	db       *db.Queries
	logger   *slog.Logger
	ghClient *GitHubClient
}

func newServer(qs *db.Queries, logger *slog.Logger) chi.Router {
	srv := &server{
		db:       qs,
		logger:   logger,
		ghClient: newGitHubClient(),
	}

	r := chi.NewRouter()
	r.Get("/github.com/{org}/{repo}/*", srv.pageShowRecipe)
	r.Get("/jobs", srv.pageListJobs)
	r.Get("/recipes", srv.pageListRecipes)

	r.Route("/api", func(r chi.Router) {
		r.Post("/github.com/{org}/{name}", srv.apiIndexRepo)
		r.Get("/metadata", srv.apiListMetadata)
	})

	return r
}
