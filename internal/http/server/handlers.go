package server

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"http-task-executor/internal/http/handlers"
)

func (s *Server) AddHandlers(router chi.Router) error {
	router.Get("/", handlers.HelloWorld())

	router.Get("/swagger/*", httpSwagger.WrapHandler)
	return nil
}
