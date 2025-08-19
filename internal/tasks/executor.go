package tasks

import "http-task-executor/internal/models"

type Executor interface {
	ExecuteTask(task models.Task)
}
