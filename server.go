package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/marcusolsson/cookhub/db/sqlc"
	"github.com/marcusolsson/cookhub/views"
	"github.com/patrickmn/go-cache"
)

type server struct {
	pool     *pgxpool.Pool
	logger   *slog.Logger
	ghClient *GitHubClient
	db       *db.Queries
	c        *cache.Cache
}

func newServer(
	pool *pgxpool.Pool,
	ghClient *GitHubClient,
	logger *slog.Logger,
	c *cache.Cache,
) chi.Router {
	srv := &server{
		pool:     pool,
		logger:   logger,
		ghClient: ghClient,
		db:       db.New(pool),
		c:        c,
	}

	r := chi.NewRouter()
	r.Get("/{provider}/{owner}/{name}/*", srv.handleError(srv.pageShowRecipe))
	r.Get("/jobs", srv.handleError(srv.pageListJobs))
	r.Get("/recipes", srv.handleError(srv.pageListRecipes))

	r.Route("/api", func(r chi.Router) {
		r.Post("/{provider}/{owner}/{name}", srv.apiIndexRepo)
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
