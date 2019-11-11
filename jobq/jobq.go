package jobq

// QManager is the interface used by other components to push to job queue
type QManager interface {
	Enqueue(url string, destination string) error

	FinishChan() <-chan []string
	FailChan() <-chan []string
}

// NewQManager returns an implementation of QManager
func NewQManager(webCapture webCapture) (QManager, error) {
	return createRabbitMQ(webCapture)
}
