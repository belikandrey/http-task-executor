package mapper

import (
	"http-task-executor/task-service/internal/task-service/models"
	"http-task-executor/task-service/internal/task-service/tasks/delivery/http/dto"
)

func MapRequestToTask(req *dto.NewTaskRequest) models.Task {
	task := models.Task{}
	task.Url = req.Url
	task.Method = req.Method
	task.Status = models.StatusNew
	task.Headers = make([]models.Header, 0)

	if len(req.Headers) > 0 {
		for name, value := range req.Headers {
			task.Headers = append(task.Headers, models.Header{Name: name, Value: value, Input: true})
		}
	}

	return task
}

func MapIdToTaskResponse(id int64) dto.NewTaskResponse {
	return dto.NewTaskResponse{Id: id}
}

func MapTaskToGetResponse(task *models.Task) dto.GetTaskResponse {
	response := dto.GetTaskResponse{ID: task.Id,
		Status:         task.Status,
		ResponseStatus: task.ResponseStatus,
		ResponseLength: task.ResponseLength}
	response.Headers = make(map[string]string)

	if len(task.Headers) > 0 {
		for _, header := range task.Headers {
			response.Headers[header.Name] = header.Value
		}
	}

	return response
}

func MapTaskToKafkaTaskMessage(task *models.Task) dto.KafkaTaskMessage {
	message := dto.KafkaTaskMessage{
		ID:     task.Id,
		Method: task.Method,
		Url:    task.Url,
	}
	message.Headers = make(map[string]string)

	if len(task.Headers) > 0 {
		for _, header := range task.Headers {
			message.Headers[header.Name] = header.Value
		}
	}

	return message
}
