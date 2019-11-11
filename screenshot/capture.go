package screenshot

import (
	"os"

	log "github.com/sirupsen/logrus"

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

type phantomJs struct {
}

func (p phantomJs) Save(url string, destination string) error {
	log.Debug("capture initiated for", url)

	cmd := exec.Command(os.Getenv("PHANTOMJS_PATH"), os.Getenv("CAPTUREJS_PATH"), url, destination)

	log.Debug("executing", cmd.String())

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("failed capturing screenshot", string(output))
	}

	return err
}
