package dto

// KafkaTaskMessage presents message received from Kafka.
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
