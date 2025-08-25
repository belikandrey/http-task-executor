package tasks

import "task-service/internal/models"

type Producer interface {
	Produce(task *models.Task) error
}
