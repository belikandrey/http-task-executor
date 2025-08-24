//go:generate mockgen -source usecase.go -destination mock/usecase.go -package mock
package tasks

import (
	"context"
	"http-task-executor/internal/models"
)

type UseCase interface {
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	GetByIdWithOutputHeaders(ctx context.Context, id int64) (*models.Task, error)
}
