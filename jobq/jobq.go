package jobq

type QManager interface {
	Enqueue(url string) error
}

func NewQManager() QManager {
	return &rabbitMQManager{}
}
