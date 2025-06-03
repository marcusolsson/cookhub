package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/marcusolsson/cookhub/db"
)

type APIServer struct {
	Store *db.Store
}

func NewRouter(store *db.Store) chi.Router {
	srv := &APIServer{
		Store: store,
	}
	r := chi.NewRouter()

	r.Post("/github.com/{org}/{name}", srv.indexRepo)
	r.Get("/github.com/{org}/{name}", srv.getIngestedFiles)
	return r
}
