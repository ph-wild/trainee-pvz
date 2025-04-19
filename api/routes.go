package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Routes(r *chi.Mux) http.Handler {
	// r := chi.NewRouter()

	r.Get("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./api/swagger.yaml")
	})

	r.Get("/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/swagger.yaml"),
	))

	return r
}
