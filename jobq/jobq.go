package jobq

type QManager interface {
	Enqueue(url string) error
}

func NewQManager(webCapture WebCapture) (QManager, error) {
	return createRabbitMQ(webCapture)
}
