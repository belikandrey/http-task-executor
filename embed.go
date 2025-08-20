package http_task_executor_v2

import "embed"

//go:embed migrations/*.sql

var MigrationFS embed.FS
