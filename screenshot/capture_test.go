package screenshot

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() {
	os.Setenv("PHANTOMJS_PATH", "/home/hamed/dev/projects/saas/phantomjs/phantomjs")
	os.Setenv("CAPTUREJS_PATH", "/home/hamed/dev/projects/saas/phantomjs/capture.js")
}
func TestCapture(t *testing.T) {
	setup()

	c, err := NewCapture()
	assert.NoError(t, err)

	err = c.Save("http://www.google.com", "/tmp/google.png")

	assert.NoError(t, err)
	assert.FileExists(t, "/tmp/google.png")
}
