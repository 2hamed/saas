package dispatcher

// Dispatcher interface is responsible for dispatching jobs to job queue and returning the result to caller
type Dispatcher interface {
	Enqueue(url string) error

	GetResult(url string) (string, error)
}

// NewDispatcher returns a new instance of Dispatcher interface
func NewDispatcher(ds dataStore, q queue) (Dispatcher, error) {
	d := &dispatcher{
		ds: ds,
		q:  q,
	}

	// TODO: start listening on finish and fail chan

	return d, nil
}

type dispatcher struct {
	ds dataStore
	q  queue
}

func (d *dispatcher) Enqueue(url string) error {
	return d.q.Enqueue(url, "")
}

func (d *dispatcher) GetResult(url string) (string, error) {

	return "", nil
}
