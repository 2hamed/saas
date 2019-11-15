package dispatcher

import (
	"errors"
	"sync"
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

		storeMutex:       &sync.Mutex{},
		setFailedMutex:   &sync.Mutex{},
		setFinishedMutex: &sync.Mutex{},
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

	assert.True(t, mockDataStore.WasStoreCalled())
	assert.True(t, mockDataStore.WasSetFinishedCalled())

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

		storeMutex:       &sync.Mutex{},
		setFailedMutex:   &sync.Mutex{},
		setFinishedMutex: &sync.Mutex{},
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

	assert.True(t, mockDataStore.WasStoreCalled())
	assert.True(t, mockDataStore.WasSetFailedCalled())
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

		storeMutex:       &sync.Mutex{},
		setFailedMutex:   &sync.Mutex{},
		setFinishedMutex: &sync.Mutex{},
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

	assert.False(t, mockDataStore.WasSetFailedCalled())
	assert.False(t, mockDataStore.WasSetFinishedCalled())

}
