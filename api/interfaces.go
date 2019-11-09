package api

type QManager interface {
	Enqueue(url string) error
}
