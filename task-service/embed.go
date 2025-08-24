package task_service

import "embed"

//go:embed migrations/*.sql

var MigrationFS embed.FS
