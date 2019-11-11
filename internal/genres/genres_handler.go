package genres

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// CreateGenresHandler -
func CreateGenresHandler(srv *Service) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	router.Post("/jointags", srv.joiningtags())
	router.Get("/list", srv.genresList(srv.stortagsPath))

	return router
}
