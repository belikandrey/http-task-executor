package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-service/internal/config"
	"task-service/internal/logger"
	"task-service/internal/tasks"
	"time"
)

const shutdownTimeout = 5 * time.Second

type Server struct {
	config   *config.Config
	database *sqlx.DB
	logger   logger.Logger
	producer tasks.Producer
}

func NewServer(config *config.Config, database *sqlx.DB, logger logger.Logger, producer tasks.Producer) *Server {
	return &Server{config: config, database: database, logger: logger, producer: producer}
}

func (s *Server) Start() error {
	s.logger.Infof("Starting server on port %d", s.config.ServerConfig.Port)

	router := chi.NewRouter()

	s.AddHandlers(router)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.ServerConfig.Host, s.config.ServerConfig.Port),
		Handler:      router,
		ReadTimeout:  s.config.ServerConfig.ReadTimeout,
		WriteTimeout: s.config.ServerConfig.WriteTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatalf("Error starting server: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	s.logger.Infof("Shutting down server on %s", sign.String())

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

	defer cancel()

	s.logger.Infof("Shutting down server properly...")
	return srv.Shutdown(ctx)
}
