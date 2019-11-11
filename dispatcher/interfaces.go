package dispatcher

type dataStore interface {
	Fetch(url string) (string, error)

	Store(url string, destination string) error
	UpdateStatus(url string, finished bool) error
	SetFailed(url string) error
}

type queue interface {
	Enqueue(url string, destination string) error

	FinishChan() <-chan []string
	FailChan() <-chan []string
}
