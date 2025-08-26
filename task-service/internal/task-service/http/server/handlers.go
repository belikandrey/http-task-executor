package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	mw "http-task-executor/task-service/internal/task-service/http/middleware"
	"http-task-executor/task-service/internal/task-service/tasks/delivery/http"
	"http-task-executor/task-service/internal/task-service/tasks/repository"
	"http-task-executor/task-service/internal/task-service/tasks/usecase"
	"time"
)

// AddHandlers added handlers to router.
func (s *Server) AddHandlers(router chi.Router) {
	s.setupMV(router)

	taskRepo := repository.NewRepository(s.database, s.logger)
	taskUseCase := usecase.NewTaskUseCase(s.logger, taskRepo, s.producer)
	taskHandlers := http.NewTaskHandlers(s.config, s.logger, taskUseCase)

	http.MapTasksRoutes(router, taskHandlers)

	router.Get("/swagger/*", httpSwagger.WrapHandler)
}

func (s *Server) setupMV(router chi.Router) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mw.New(s.logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.URLFormat)
}
