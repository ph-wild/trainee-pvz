package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"trainee-pvz/internal/openapi"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetPvz(w http.ResponseWriter, r *http.Request, pvz openapi.GetPvzParams) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) PostDummyLogin(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) PostLogin(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) PostProducts(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) PostPvz(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) PostPvzPvzIdCloseLastReception(w http.ResponseWriter, r *http.Request, uid uuid.UUID) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) PostPvzPvzIdDeleteLastProduct(w http.ResponseWriter, r *http.Request, uid uuid.UUID) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) PostReceptions(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) PostRegister(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) Routes() *chi.Mux {
	router := chi.NewRouter()

	return router
}
