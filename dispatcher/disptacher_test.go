package dispatcher

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnqueueSuccess(t *testing.T) {
	mockDataStore := &mockDataStore{
		fetchResult: "",
		fetchErr:    nil,

		fetchStatusExists:     false,
		fetchStatusIsFinished: false,
		fetchStatusErr:        nil,

		storeErr: nil,

		setFinishedErr: nil,
		setFailedErr:   nil,
	}

	mockQueue := &mockQueue{
		failJob:    false,
		finishChan: make(chan []string),
		failChan:   make(chan []string),
	}

	d, err := NewDispatcher(mockDataStore, mockQueue)
	assert.NoError(t, err)

	err = d.Enqueue("some url")
	assert.NoError(t, err)

	assert.True(t, mockDataStore.wasStoreCalled)
	assert.True(t, mockDataStore.wasSetFinishedCalled)

}

func TestEnqueueJobFail(t *testing.T) {
	mockDataStore := &mockDataStore{
		fetchResult: "",
		fetchErr:    nil,

		fetchStatusExists:     false,
		fetchStatusIsFinished: false,
		fetchStatusErr:        nil,

		storeErr: nil,

		setFinishedErr: nil,
		setFailedErr:   nil,
	}

	mockQueue := &mockQueue{
		failJob:    true,
		finishChan: make(chan []string),
		failChan:   make(chan []string),
	}

	d, err := NewDispatcher(mockDataStore, mockQueue)
	assert.NoError(t, err)

	err = d.Enqueue("some url")
	assert.NoError(t, err)

	assert.True(t, mockDataStore.wasStoreCalled)
	assert.True(t, mockDataStore.wasSetFailedCalled)
}

func TestEnqueueStoreFailed(t *testing.T) {
	mockDataStore := &mockDataStore{
		fetchResult: "",
		fetchErr:    nil,

		fetchStatusExists:     false,
		fetchStatusIsFinished: false,
		fetchStatusErr:        nil,

		storeErr: errors.New("arbitrary error"),

		setFinishedErr: nil,
		setFailedErr:   nil,
	}

	mockQueue := &mockQueue{
		failJob:    false,
		finishChan: make(chan []string),
		failChan:   make(chan []string),
	}

	d, err := NewDispatcher(mockDataStore, mockQueue)
	assert.NoError(t, err)

	err = d.Enqueue("some url")
	assert.Error(t, err)

	assert.False(t, mockDataStore.wasSetFailedCalled)
	assert.False(t, mockDataStore.wasSetFinishedCalled)

}
