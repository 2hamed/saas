package storage

// DataStore is the abstract type to work with the underlying datastore
type DataStore interface {
	Fetch(url string) (string, error)
	FetchStatus(url string) (exists bool, isPending bool, isFinished bool, err error)

	Store(url string, destination string) error
	UpdatePath(url string, destination string) error
	SetFinished(url string) error
	SetFailed(url string) error
	Delete(url string) error
}

// NewDataStore returns an implementation of DataStore
func NewDataStore() (DataStore, error) {
	return createMongoDataStore()
}
