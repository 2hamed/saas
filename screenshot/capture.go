package screenshot

// Capture is the interface to capture the screenshot
type Capture interface {
	Save(url string) (path string, err error)
}

// NewCapture returns a concrete implementaion of the Capture interface
func NewCapture() Capture {
	return phantomJs{}
}

type phantomJs struct {
}

func (p phantomJs) Save(url string) (path string, err error) {
	// os.Exec("")
	return "", nil
}
