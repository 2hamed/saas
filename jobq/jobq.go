package jobq

import "fmt"

// QManager is the interface used by other components to push to job queue
type QManager interface {
	Enqueue(url string, destination string) error

	FinishChan() <-chan []string
	FailChan() <-chan []string
}

// NewQManager returns an implementation of QManager
func NewQManager(webCapture webCapture) (QManager, error) {
	conn, err := createRabbitMQConnection()
	if err != nil {
		return nil, fmt.Errorf("failed connecting to rabbitmq instance: %v", err)
	}

	return createRabbitMQ(webCapture, conn)
}
