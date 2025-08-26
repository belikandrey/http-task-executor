package taskservice

import "embed"

//go:embed migrations/*.sql

// MigrationFS - needed to execute db migrations.
var MigrationFS embed.FS
