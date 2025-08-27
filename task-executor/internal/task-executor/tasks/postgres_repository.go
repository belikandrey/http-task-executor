//go:generate mockgen -source postgres_repository.go -destination mock/postgres_repository.go -package mock
package tasks

import (
	"context"

	"http-task-executor/task-executor/internal/task-executor/models"
)

// Repository represents db repository to work with models.Task.
type Repository interface {
	// UpdateResult updates task result fields
	UpdateResult(ctx context.Context, task *models.Task) error
	// UpdateStatus updates task status field
	UpdateStatus(ctx context.Context, id int64, newStatus string) error
}
