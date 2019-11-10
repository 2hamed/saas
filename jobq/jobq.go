package jobq

// QManager is the interface used by other components to push to job queue
type QManager interface {
	Enqueue(url string) error
}

// NewQManager returns an implementation of QManager
func NewQManager(webCapture webCapture, jc jobCallback) (QManager, error) {
	return createRabbitMQ(webCapture, jc)
}
