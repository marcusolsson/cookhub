package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
	"github.com/marcusolsson/cookhub/views"
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
	r.Get("/github.com/{org}/{repo}/*", srv.handleError(srv.pageShowRecipe))
	r.Get("/jobs", srv.handleError(srv.pageListJobs))
	r.Get("/recipes", srv.handleError(srv.pageListRecipes))

	r.Route("/api", func(r chi.Router) {
		r.Post("/github.com/{org}/{name}", srv.apiIndexRepo)
		r.Get("/metadata", srv.apiListMetadata)
	})

	return r
}

func (s *server) handleError(
	handler func(w http.ResponseWriter, req *http.Request) error,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := handler(w, req); err != nil {
			s.logger.Error(err.Error())

			if err := views.ErrorPage().Render(req.Context(), w); err != nil {
				s.logger.Error(err.Error())
			}
		}
	}
}
