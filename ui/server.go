package ui

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
	"github.com/marcusolsson/cookhub/pkg/recipe"
)

type server struct {
	db *db.Queries
}

func NewRouter(qs *db.Queries) chi.Router {
	srv := &server{
		db: qs,
	}
	r := chi.NewRouter()
	r.Get("/github.com/{org}/{name}", srv.getRecipe)
	return r
}

func (s *server) getRecipe(w http.ResponseWriter, req *http.Request) {
	var (
		org  = chi.URLParam(req, "org")
		name = chi.URLParam(req, "name")
		slug = fmt.Sprintf("%s/%s", org, name)
	)

	ctx := req.Context()

	jobID, err := s.db.GetLatestJobBySlug(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	files, err := s.db.GetFilesByJob(ctx, jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	file := files[3]

	cooklangRecipe, err := recipe.ParseCooklangRecipe(file.Name, file.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	schemaOrgRecipe := recipe.ConvertCooklangToSchemaOrg(file.Name, cooklangRecipe)

	component := recipeView(schemaOrgRecipe)
	component.Render(req.Context(), w)
}
