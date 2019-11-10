package screenshot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() {
	phantomJsPath = "/home/hamed/dev/projects/saas/phantomjs/phantomjs"
	captureJsPath = "/home/hamed/dev/projects/saas/phantomjs/capture.js"
}
func TestCapture(t *testing.T) {
	setup()

	c, err := NewCapture()
	assert.NoError(t, err)

	err = c.Save("http://www.google.com", "/tmp/google.png")

	assert.NoError(t, err)
	assert.FileExists(t, "/tmp/google.png")
}
