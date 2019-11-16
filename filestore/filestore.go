package filestore

// FileStore is the interface to store files on an abstract block storage
type FileStore interface {
	Store(path string) (string, error)
}

// NewFileStore returns a concrete implementation of FileStore
func NewFileStore() (FileStore, error) {
	return createMinio()
}
