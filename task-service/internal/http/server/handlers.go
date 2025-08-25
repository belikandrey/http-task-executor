package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	mw "task-service/internal/http/middleware"
	taskHttp "task-service/internal/tasks/delivery/http"
	"task-service/internal/tasks/repository"
	"task-service/internal/tasks/usecase"
	"time"
)

func (s *Server) AddHandlers(router chi.Router) {
	s.setupMV(router)

	taskRepo := repository.NewRepository(s.database, s.logger)
	//taskExec := executor.NewExecutor(s.logger, taskRepo, &executor.ClientProvider{}, s.config.ExternalServiceTimeout)
	taskUseCase := usecase.NewTaskUseCase(s.logger, taskRepo, s.producer)
	taskHandlers := taskHttp.NewTaskHandlers(s.config, s.logger, taskUseCase)

	taskHttp.MapTasksRoutes(router, taskHandlers)

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
