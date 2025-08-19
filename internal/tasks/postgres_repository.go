package tasks

import (
	"context"
	"http-task-executor/internal/models"
)

type Repository interface {
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	GetByIdWithOutputHeaders(ctx context.Context, id int64) (*models.Task, error)
	UpdateStatus(ctx context.Context, id int64, newStatus string) error
	UpdateResult(ctx context.Context, task *models.Task) error
}
