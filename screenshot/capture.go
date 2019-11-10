package screenshot

import (
	"log"
	"os"
	"os/exec"
)

// Capture is the interface to capture the screenshot
type Capture interface {
	Save(url string, destination string) error
}

// NewCapture returns a concrete implementaion of the Capture interface
func NewCapture() (Capture, error) {
	return phantomJs{}, nil
}

var (
	phantomJsPath = os.Getenv("PHANTOMJS_PATH")
	captureJsPath = os.Getenv("CAPTUREJS_PATH")
)

type phantomJs struct {
}

func (p phantomJs) Save(url string, destination string) error {
	cmd := exec.Command(phantomJsPath, captureJsPath, url, destination)
	output, err := cmd.Output()

	if err != nil {
		log.Println(string(output))
	}
	return err
}
