package artists

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// CreateArtistsHandler -
func CreateArtistsHandler(srv *Service) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	router.Get("/", srv.artistsList())
	return router
}
