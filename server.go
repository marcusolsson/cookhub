package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/marcusolsson/cookhub/db/sqlc"
	"github.com/marcusolsson/cookhub/github"
	"github.com/marcusolsson/cookhub/views"
	"github.com/patrickmn/go-cache"
)

type Server struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
	gh     *github.Client
	db     *db.Queries
	c      *cache.Cache
}

// Router returns a new chi router with the server's routes configured.
func (srv *Server) Router() chi.Router {
	r := chi.NewRouter()

	h := errorHandler(srv.logger)

	r.Use(middleware.StripSlashes)

	r.Get("/{provider}/{owner}/{repo_name}/*", h(srv.recipePage))
	r.Get("/{provider}/{owner}/{repo_name}", h(srv.cookbookPage))

	// API routes handle admin operations.
	r.Route("/api", func(r chi.Router) {
		r.Post("/repos", srv.apiAddRepository)
		r.Post("/import", srv.apiImportRepos)
	})

	// Serve all files in the "/static" directory.
	r.Handle("/static/*", http.StripPrefix("/static/", StaticFileServer))
	r.Handle("/", http.RedirectHandler("/github.com/marcusolsson/recipes", http.StatusFound))

	r.NotFound(h(srv.notFoundPage))

	return r
}

// errorHandler is an adapter to allow handlers to return errors,
type errHandlerFunc func(w http.ResponseWriter, req *http.Request) error

// errorHandler is a middleware that wraps an error handler function.
func errorHandler(logger *slog.Logger) func(errHandlerFunc) http.HandlerFunc {
	return func(handler errHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			if err := handler(w, req); err != nil {
				logger.Error(err.Error())
				if err := views.ErrorPage().Render(req.Context(), w); err != nil {
					logger.Error(err.Error())
				}
			}
		}
	}
}
