package storage

type DataFetcher interface {
	Fetch(url string) (string, error)
}

type DataStore interface {
	Store(url string, path string) error
}
