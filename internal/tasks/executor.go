//go:generate mockgen -source executor.go -destination mock/executor.go -package mock

package tasks

import "http-task-executor/internal/models"

type Executor interface {
	ExecuteTask(task models.Task)
}
