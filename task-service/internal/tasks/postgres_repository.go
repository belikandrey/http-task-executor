//go:generate mockgen -source postgres_repository.go -destination mock/postgres_repository.go -package mock
package tasks

import (
	"context"
	"task-service/internal/models"
)

type Repository interface {
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	GetByIdWithOutputHeaders(ctx context.Context, id int64) (*models.Task, error)
	UpdateStatus(ctx context.Context, id int64, newStatus string) error
	Delete(ctx context.Context, id int64) error
}
