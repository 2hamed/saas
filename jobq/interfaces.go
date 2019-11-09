package jobq

type WebCapture interface {
	Save(url string, destination string) error
}
