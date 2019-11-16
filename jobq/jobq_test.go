package jobq

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockCapture struct {
	saveErr error
}

func (m *mockCapture) Save(url string, destination string) error {
	return m.saveErr
}
func TestJobQFinish(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	jobq, err := NewQManager(&mockCapture{})
	assert.NoError(t, err)

	err = jobq.Enqueue("url", "path")
	assert.NoError(t, err)

	timeout := make(chan struct{})

	go func() {
		time.Sleep(10 * time.Second)
		timeout <- struct{}{}
	}()

	select {
	case <-timeout:
		assert.Fail(t, "test timed out")
	case res := <-jobq.FinishChan():
		assert.Equal(t, res, []string{"url", "path"})
	case <-jobq.FailChan():
		assert.Fail(t, "should not have issued anything to failChan")
	}

	jobq.CleanUp()
}

func TestJobQFail(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	jobq, err := NewQManager(&mockCapture{
		saveErr: errors.New("arbitrary error"),
	})

	assert.NoError(t, err)

	err = jobq.Enqueue("url", "path")
	assert.NoError(t, err)

	timeout := make(chan struct{})

	go func() {
		time.Sleep(10 * time.Second)
		timeout <- struct{}{}
	}()

	select {
	case <-timeout:
		assert.Fail(t, "test timed out")
	case <-jobq.FinishChan():
		assert.Fail(t, "should not have issued anything to finishChan")
	case res := <-jobq.FailChan():
		assert.Equal(t, res, []string{"url", "path"})
	}
	jobq.CleanUp()
}
