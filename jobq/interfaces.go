package jobq

type webCapture interface {
	Save(url string, destination string) error
}
