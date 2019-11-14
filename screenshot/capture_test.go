package screenshot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapture(t *testing.T) {
	if testing.Short() {
		t.Skip("test skipped in short mode")
	}

	c, err := NewCapture()
	assert.NoError(t, err)

	err = c.Save("http://www.google.com", "/tmp/google.png")

	assert.NoError(t, err)
	assert.FileExists(t, "/tmp/google.png")
}
