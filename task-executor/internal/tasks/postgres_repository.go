//go:generate mockgen -source postgres_repository.go -destination mock/postgres_repository.go -package mock
package tasks

import (
	"context"
	"task-executor/internal/models"
)

type Repository interface {
	UpdateResult(ctx context.Context, task *models.Task) error
	UpdateStatus(ctx context.Context, id int64, newStatus string) error
}
