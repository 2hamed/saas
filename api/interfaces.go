package api

type dispatcher interface {
	Enqueue(url string) error
	GetResult(url string) (string, error)
	GetStatus(url string) (exists bool, isFinished bool, err error)
}
