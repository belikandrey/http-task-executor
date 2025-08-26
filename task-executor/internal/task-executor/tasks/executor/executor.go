package executor

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"task-executor/internal/task-executor/dto"
	"task-executor/internal/task-executor/logger"
	"task-executor/internal/task-executor/models"
	tasks2 "task-executor/internal/task-executor/tasks"
	"time"
)

type Executor struct {
	log            logger.Logger
	repo           tasks2.Repository
	timeout        time.Duration
	clientProvider tasks2.ClientProvider
}

type ClientProvider struct {
	Timeout time.Duration
}

func (c *ClientProvider) Client() *http.Client {
	return &http.Client{Timeout: c.Timeout}
}

func NewExecutor(log logger.Logger, repo tasks2.Repository, clientProvider tasks2.ClientProvider, timeout time.Duration) *Executor {
	return &Executor{log: log, repo: repo, clientProvider: clientProvider, timeout: timeout}
}

func (e *Executor) Execute(value []byte) {
	e.log.Debugf("Executing message: %v", string(value))
	var message dto.KafkaTaskMessage

	err := json.Unmarshal(value, &message)

	if err != nil {
		e.log.Errorf("executor.ExecuteTask Error unmarshal message : %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)

	defer cancel()

	err = e.repo.UpdateStatus(ctx, message.ID, models.StatusInProcess)
	if err != nil {
		e.log.Errorf("executor.ExecuteTask.UpdateStatus : %v", err)
		return
	}
	e.log.Infof("executor.ExecuteTask: updatedTask %v with method %s and url %s", message.ID, message.Method, message.Url)
	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(message.Method), message.Url, nil)
	if err != nil {
		e.setErrorStatus(message.ID)
		e.log.Errorf("executor.ExecuteTask.NewRequestWithContext : %v", err)
		return
	}
	if message.Headers != nil {
		for key, val := range message.Headers {
			req.Header.Add(key, val)
		}
	}

	client := e.clientProvider.Client()

	resp, err := client.Do(req)
	if err != nil {
		e.setErrorStatus(message.ID)
		e.log.Errorf("executor.ExecuteTask.DoRequest : %v", err)
		return
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			e.log.Errorf("executor.ExecuteTask.DoBody.Close : %v", err)
		}
	}()

	contentLength, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		e.setErrorStatus(message.ID)
		e.log.Errorf("executor.ExecuteTask.DoRequest.Copy : %v", err)
		return
	}

	e.log.Infof("executor.ExecuteTask: updatedTask %v with method %s and url %s executed successfully with code %v", message.ID, message.Method, req.URL, resp.StatusCode)

	updatedTask := models.Task{Id: message.ID}

	updatedTask.ResponseLength = &contentLength
	updatedTask.Status = models.StatusDone
	code := int64(resp.StatusCode)
	updatedTask.ResponseStatus = &code
	outputHeaders := make([]models.Header, 0)
	for k, v := range resp.Header {
		res := strings.Join(v[:], ",")
		outputHeaders = append(outputHeaders, models.Header{Name: k, Value: res, Input: false})
	}

	updatedTask.Headers = append(updatedTask.Headers, outputHeaders...)
	err = e.repo.UpdateResult(ctx, &updatedTask)
	if err != nil {
		e.setErrorStatus(updatedTask.Id)
		e.log.Errorf("executor.ExecuteTask.UpdateResult : %v", err)
	}
}

func (e *Executor) setErrorStatus(id int64) {
	err := e.repo.UpdateStatus(context.Background(), id, models.StatusError)
	if err != nil {
		e.log.Errorf("executor.ExecuteTask.setErrorStatus.UpdateStatus : %v", err)
	}
}
