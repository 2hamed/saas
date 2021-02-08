package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapture(t *testing.T) {

	t.Skip("this test should be run on a host which supports PhantomJS")

	c, err := NewCapture()
	assert.NoError(t, err)

	err = c.Save("http://www.google.com", "google.png")

	assert.NoError(t, err)
	assert.FileExists(t, "google.png")
}
