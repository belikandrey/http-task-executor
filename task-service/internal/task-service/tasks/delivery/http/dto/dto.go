package dto

// NewTaskRequest represents request to create new task.
// It contains URL, Method, Headers.
type NewTaskRequest struct {
	// URL - url to send request.
	URL string `json:"url"`
	// Method - http method to use in request.
	Method string `json:"method"`
	// Headers - map of headers that will be used in request.
	Headers map[string]string `json:"headers"`
}

// NewTaskResponse represents response to create new task request.
// It contains task ID.
type NewTaskResponse struct {
	// ID - unique identifier of task.
	ID int64 `json:"id"`
}

// GetTaskResponse represents response to get task request.
// It contains task ID, Status, ResponseStatus, ResponseLength, Headers.
type GetTaskResponse struct {
	// ID - unique identifier of task.
	ID int64 `json:"id"`
	// Status - current task status.
	Status string `json:"status"`
	// ResponseStatus - response status from the 3-rd service.
	ResponseStatus *int64 `json:"httpStatusCode"`
	// ResponseLength - response length from the 3-rd service.
	ResponseLength *int64 `json:"length"`
	// Headers - response headers from the 3-rd service.
	Headers map[string]string `json:"headers"`
}

// KafkaTaskMessage represents message received from Kafka.
// It contains ID, Method, URL, Headers.
type KafkaTaskMessage struct {
	// ID - unique identifier of task.
	ID int64 `json:"id"`
	// Method - http method to use in request.
	Method string `json:"method"`
	// URL - url to send request.
	URL string `json:"url"`
	// Headers - map of headers that will be used in request.
	Headers map[string]string `json:"headers"`
}
