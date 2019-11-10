package storage

type DataFetcher interface {
	Fetch(url string) (string, error)
}

type DataStore interface {
	Store(url string, destination string) error
	UpdateStatus(url string, finished bool) error
	SetFailed(url string) error
}
