package dispatcher

type dataStore interface {
	Fetch(url string) (string, error)
	FetchStatus(url string) (exists bool, isPending bool, isFinished bool, err error)

	Store(url string, destination string) error
	SetFinished(url string) error
	SetFailed(url string) error
}

type queue interface {
	Enqueue(url string, destination string) error
	FinishChan() <-chan []string
	FailChan() <-chan []string
}
