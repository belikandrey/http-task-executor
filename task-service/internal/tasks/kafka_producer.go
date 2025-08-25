//go:generate mockgen -source kafka_producer.go -destination mock/kafka_producer.go -package mock

package tasks

import "task-service/internal/models"

type Producer interface {
	Produce(task *models.Task) error
}
