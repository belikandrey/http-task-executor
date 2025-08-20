//go:generate mockgen -source executor.go -destination mock/executor.go -package mock

package tasks

import (
	"http-task-executor/internal/models"
	"net/http"
)

type Executor interface {
	ExecuteTask(task models.Task)
}

type ClientProvider interface {
	Client() *http.Client
}
