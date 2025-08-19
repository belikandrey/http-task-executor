package executor

import (
	"context"
	"http-task-executor/internal/logger"
	"http-task-executor/internal/models"
	"http-task-executor/internal/tasks"
	"io"
	"net/http"
	"strings"
	"time"
)

type Executor struct {
	log     logger.Logger
	repo    tasks.Repository
	timeout time.Duration
}

func NewExecutor(log logger.Logger, repo tasks.Repository, timeout time.Duration) *Executor {
	return &Executor{log: log, repo: repo, timeout: timeout}
}

func (e *Executor) ExecuteTask(task models.Task) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)

	defer cancel()

	err := e.repo.UpdateStatus(ctx, task.Id, models.StatusInProcess)
	if err != nil {
		e.log.Errorf("executor.ExecuteTask.UpdateStatus : %v", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, task.Method, task.Url, nil)
	if err != nil {
		e.setErrorStatus(task.Id)
		e.log.Errorf("executor.ExecuteTask.NewRequestWithContext : %v", err)
		return
	}
	if task.Headers != nil {
		for _, v := range task.Headers {
			req.Header.Add(v.Name, v.Value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		e.setErrorStatus(task.Id)
		e.log.Errorf("executor.ExecuteTask.DoRequest : %v", err)
		return
	}

	defer resp.Body.Close()

	contentLength, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		e.setErrorStatus(task.Id)
		e.log.Errorf("executor.ExecuteTask.DoRequest.Copy : %v", err)
		return
	}

	task.ResponseLength = &contentLength
	task.Status = models.StatusDone
	code := int64(resp.StatusCode)
	task.ResponseStatus = &code
	outputHeaders := make([]models.Header, 0)
	for k, v := range resp.Header {
		res := strings.Join(v[:], ",")
		outputHeaders = append(outputHeaders, models.Header{Name: k, Value: res, Input: false})
	}

	task.Headers = append(task.Headers, outputHeaders...)
	err = e.repo.UpdateResult(ctx, &task)
	if err != nil {
		e.setErrorStatus(task.Id)
		e.log.Errorf("executor.ExecuteTask.UpdateResult : %v", err)
	}
}

func (e *Executor) setErrorStatus(id int64) {
	err := e.repo.UpdateStatus(context.Background(), id, models.StatusError)
	if err != nil {
		e.log.Errorf("executor.ExecuteTask.setErrorStatus.UpdateStatus : %v", err)
	}
}
