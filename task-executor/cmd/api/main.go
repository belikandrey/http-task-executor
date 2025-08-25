package main

import (
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"os/signal"
	"syscall"
	"task-executor/internal/config"
	"task-executor/internal/logger"
	"task-executor/internal/postgres"
	"task-executor/internal/tasks/consumer"
	"task-executor/internal/tasks/executor"
	"task-executor/internal/tasks/repository"
)

func main() {
	log.Println("Starting task executor")

	appConfig := config.MustLoad()

	appLogger, err := logger.NewLogger(appConfig)

	if err != nil {
		log.Fatalf("Init logger error: %v", err)
	}

	appLogger.Infof("Env: %s, LogLevel: %s", appConfig.Env, appConfig.LoggerConfig.Level)

	database, err := postgres.NewPostgresqlDatabase(appConfig)
	if err != nil {
		appLogger.Fatalf("Init postgresql database error: %v", err)
	} else {
		appLogger.Infof("Init postgresql database success")
	}

	defer func(database *sqlx.DB) {
		err := database.Close()
		if err != nil {
			appLogger.Errorf("Close postgresql database error: %v", err)
		}
	}(database)

	repo := repository.NewRepository(database, appLogger)

	clientProvider := executor.ClientProvider{}

	exec := executor.NewExecutor(appLogger, repo, &clientProvider, appConfig.ExternalServiceTimeout)

	cons, err := consumer.NewKafkaConsumer(appConfig, exec, appLogger)
	if err != nil {
		appLogger.Fatalf("Init consumer error: %v", err)
	} else {
		appLogger.Info("Init consumer success")
	}
	go func() {
		cons.Start()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	appLogger.Infof("Shutting down server on %s", sign.String())

	appLogger.Infof("Shutting down server properly...")

	err = cons.Close()
	if err != nil {
		appLogger.Fatalf("Close consumer error: %v", err)
	}
}
