package mapper

import (
	"http-task-executor/task-service/internal/task-service/models"
	"http-task-executor/task-service/internal/task-service/tasks/delivery/http/dto"
)

const (
	emptyMapLen = 0
)

// MapRequestToTask - maps dto.NewTaskRequest to models.Task.
func MapRequestToTask(req *dto.NewTaskRequest) models.Task {
	task := models.Task{}
	task.URL = req.URL
	task.Method = req.Method
	task.Status = models.StatusNew
	task.Headers = make([]models.Header, emptyMapLen)

	if len(req.Headers) > emptyMapLen {
		for name, value := range req.Headers {
			task.Headers = append(
				task.Headers,
				models.Header{Name: name, Value: value, Input: true},
			)
		}
	}

	return task
}

// MapIDToTaskResponse - maps ID to dto.NewTaskResponse.
func MapIDToTaskResponse(id int64) dto.NewTaskResponse {
	return dto.NewTaskResponse{ID: id}
}

// MapTaskToGetResponse - maps models.Task to dto.GetTaskResponse.
func MapTaskToGetResponse(task *models.Task) dto.GetTaskResponse {
	response := dto.GetTaskResponse{
		ID:             task.ID,
		Status:         task.Status,
		ResponseStatus: task.ResponseStatus,
		ResponseLength: task.ResponseLength,
	}
	response.Headers = make(map[string]string)

	if len(task.Headers) > emptyMapLen {
		for _, header := range task.Headers {
			response.Headers[header.Name] = header.Value
		}
	}

	return response
}

// MapTaskToKafkaTaskMessage - maps models.Task to dto.KafkaTaskMessage.
func MapTaskToKafkaTaskMessage(task *models.Task) dto.KafkaTaskMessage {
	message := dto.KafkaTaskMessage{
		ID:     task.ID,
		Method: task.Method,
		URL:    task.URL,
	}
	message.Headers = make(map[string]string)

	if len(task.Headers) > emptyMapLen {
		for _, header := range task.Headers {
			message.Headers[header.Name] = header.Value
		}
	}

	return message
}
