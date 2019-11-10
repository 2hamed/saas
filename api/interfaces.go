package api

type coordinator interface {
	CaptureAsync(url string) error
}
