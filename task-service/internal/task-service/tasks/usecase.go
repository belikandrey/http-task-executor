//go:generate mockgen -source usecase.go -destination mock/usecase.go -package mock
package tasks

import (
	"context"

	"http-task-executor/task-service/internal/task-service/models"
)

// UseCase represents service layer to models.Task.
type UseCase interface {
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	GetByIDWithOutputHeaders(ctx context.Context, id int64) (*models.Task, error)
}
