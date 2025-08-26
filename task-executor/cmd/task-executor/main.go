package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"task-executor/internal/task-executor/config"
	"task-executor/internal/task-executor/logger"
	"task-executor/internal/task-executor/postgres"
	"task-executor/internal/task-executor/tasks/consumer"
	"task-executor/internal/task-executor/tasks/executor"
	"task-executor/internal/task-executor/tasks/repository"
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
	}
	appLogger.Infof("Init postgresql database success")

	defer func() {
		err := database.Close()
		if err != nil {
			appLogger.Errorf("Close postgresql database error: %v", err)
		}
	}()

	repo := repository.NewRepository(database, appLogger)

	clientProvider := executor.ClientProvider{Timeout: appConfig.ExternalServiceTimeout}

	exec := executor.NewExecutor(appLogger, repo, &clientProvider, appConfig.ExternalServiceTimeout)

	cons, err := consumer.NewKafkaConsumer(appConfig, exec, appLogger)
	if err != nil {
		appLogger.Fatalf("Init consumer error: %v", err)
	}
	appLogger.Info("Init consumer success")

	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		cons.Start(ctx)
	}()
	defer cancelFunc()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	appLogger.Infof("Shutting down server on %s", sign.String())

	appLogger.Infof("Shutting down server properly...")

	err = cons.Close(cancelFunc)
	if err != nil {
		appLogger.Fatalf("Close consumer error: %v", err)
	}
}
