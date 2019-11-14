package api

type dispatcher interface {
	Enqueue(url string) error
	FetchResult(url string) (string, error)
	FetchStatus(url string) (exists bool, isPending bool, isFinished bool, err error)
}
