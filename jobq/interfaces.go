package jobq

type webCapture interface {
	Save(url string, destination string) error
}

type jobCallback interface {
	JobFailed(url string)
	JobFinished(url string)
}
