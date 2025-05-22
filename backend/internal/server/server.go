package server

import (
	"log"
	"net/http"

	"time"

	"github.com/Chaitu35/claimeasy/backend/ent"
	auth "github.com/Chaitu35/claimeasy/backend/internal/auth/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	client *ent.Client	
	Router *chi.Mux
}

func New(client *ent.Client) *Server {
	srv := &Server{
		client: client,
		Router: chi.NewRouter(),
	}
	srv.Router.Use(middleware.Logger)
	srv.Router.Use(middleware.Recoverer)
	srv.Router.Use(middleware.Timeout(60 * time.Second))
	srv.Router.Use(middleware.RequestID)
	srv.Router.Use(middleware.RealIP)
	srv.Router.Use(middleware.StripSlashes)
	srv.Router.Use(middleware.Compress(5))
	srv.routes()
	return srv
}

func (s *Server) Start(addr string){
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, s.Router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Server) routes() {
	s.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})	
	s.Router.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Version 1.0.0"))
	})

	s.Router.Mount("/auth", auth.Routes(s.client))

}
