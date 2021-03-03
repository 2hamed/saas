package odin

import (
	"context"
	"fmt"
)

type CaptureJob struct {
	UUID string `json:"uuid"`
	URL  string `json:"url"`

	Ack  func() `json:"-"`
	Nack func() `json:"-"`
}

// QManager is the interface used by other components to push to job queue
type QManager interface {
	Enqueue(ctx context.Context, job CaptureJob) error

	GetJobChan(ctx context.Context) (<-chan CaptureJob, error)

	CleanUp()
}

// NewQManager returns an implementation of QManager
func NewQManager() (QManager, error) {
	conn, err := createRabbitMQConnection()
	if err != nil {
		return nil, fmt.Errorf("failed connecting to rabbitmq instance: %w", err)
	}

	return createRabbitMQ(conn)
}
