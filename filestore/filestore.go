package filestore

type FileStore interface {
	Store(path string) (string, error)
}

func NewFileStore() (FileStore, error) {
	return createMinio()
}
