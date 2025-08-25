package main

import (
	"github.com/jmoiron/sqlx"
	"log"
	_ "task-service/docs"
	"task-service/internal/config"
	"task-service/internal/http/server"
	"task-service/internal/logger"
	"task-service/internal/migration"
	"task-service/internal/postgres"
	"task-service/internal/tasks/producer"
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

	defer func(database *sqlx.DB) {
		err := database.Close()
		if err != nil {
			appLogger.Errorf("Close postgresql database error: %v", err)
		}
	}(database)

	err = migration.MigratePostgresql(database)
	if err != nil {
		appLogger.Fatalf("MigratePostgresql database error: %v", err)
	} else {
		appLogger.Infof("Database migrated successfully")
	}

	produce, err := producer.NewTaskProducer(appConfig, appLogger)
	if err != nil {
		appLogger.Fatalf("Init kafka producer error: %v", err)
	}
	defer func() {
		produce.Close()
	}()

	httpServer := server.NewServer(appConfig, database, appLogger, produce)

	if err := httpServer.Start(); err != nil {
		appLogger.Fatal("Start server error: %v", err)
	}
}
