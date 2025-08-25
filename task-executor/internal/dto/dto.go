package dto

type KafkaTaskMessage struct {
	ID      int64             `json:"id"`
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}
