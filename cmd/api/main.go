package main

import (
	_ "http-task-executor/docs"
	"http-task-executor/internal/config"
	"http-task-executor/internal/http/server"
	"http-task-executor/internal/logger"
	"http-task-executor/internal/postgres"
	"log"
)

// @title Task executor Rest API
// @version 1.0
// @description Executes request to 3-rd services
// @contact.name Andrei Belik
// @contact.url https://github.com/belikandrey
// @contact.email belikandrey01@gmail.com
// @BasePath /
func main() {
	log.Println("Starting api server")

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

	defer database.Close()

	httpServer := server.NewServer(appConfig, database, appLogger)

	if err := httpServer.Start(); err != nil {
		appLogger.Fatal("Start server error: %v", err)
	}
}
