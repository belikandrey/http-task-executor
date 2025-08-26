//go:generate mockgen -source kafka_producer.go -destination mock/kafka_producer.go -package mock

package tasks

import (
	"http-task-executor/task-service/internal/task-service/models"
)

// Producer presents message producer to Kafka.
type Producer interface {
	Produce(task *models.Task) error
}
