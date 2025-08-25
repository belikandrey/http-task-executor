package dto

type NewTaskRequest struct {
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
}

type NewTaskResponse struct {
	Id int64 `json:"id"`
}

type GetTaskResponse struct {
	ID             int64             `json:"id"`
	Status         string            `json:"status"`
	ResponseStatus *int64            `json:"httpStatusCode"`
	ResponseLength *int64            `json:"length"`
	Headers        map[string]string `json:"headers"`
}

type KafkaTaskMessage struct {
	ID      int64             `json:"id"`
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}
