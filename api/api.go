package api

import (
	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
)

type APIServer struct {
	db *db.Queries
}

func NewRouter(db *db.Queries) chi.Router {
	srv := &APIServer{
		db: db,
	}
	r := chi.NewRouter()

	r.Post("/github.com/{org}/{name}", srv.indexRepo)
	r.Get("/github.com/{org}/{name}", srv.getIngestedFiles)
	return r
}
