package dispatcher

import (
	"sync"
	"time"
)

type mockDataStore struct {
	fetchResult string
	fetchErr    error

	fetchStatusExists     bool
	fetchStatusIsFinished bool
	fetchStatusIsPending  bool
	fetchStatusErr        error

	storeErr error

	setFinishedErr error
	setFailedErr   error

	wasFetchCalled       bool
	wasFetchStatusCalled bool
	wasStoreCalled       bool
	wasSetFailedCalled   bool
	wasSetFinishedCalled bool

	// These mutexes are just to avoid race conditions in tests
	setFailedMutex   *sync.Mutex
	setFinishedMutex *sync.Mutex
	storeMutex       *sync.Mutex
}

func (m *mockDataStore) Fetch(url string) (string, error) {
	m.wasFetchCalled = true
	return m.fetchResult, m.fetchErr
}

func (m *mockDataStore) FetchStatus(url string) (exists bool, isPending bool, isFinished bool, err error) {
	m.wasFetchStatusCalled = true
	return m.fetchStatusExists, m.fetchStatusIsPending, m.fetchStatusIsFinished, m.fetchStatusErr
}

func (m *mockDataStore) Store(url string, destination string) error {
	m.storeMutex.Lock()
	m.wasStoreCalled = true
	m.storeMutex.Unlock()

	return m.storeErr
}

func (m *mockDataStore) SetFinished(url string) error {
	m.setFinishedMutex.Lock()
	m.wasSetFinishedCalled = true
	m.setFinishedMutex.Unlock()

	return m.setFinishedErr
}

func (m *mockDataStore) SetFailed(url string) error {
	m.setFailedMutex.Lock()
	m.wasSetFailedCalled = true
	m.setFailedMutex.Unlock()

	return m.setFailedErr
}

func (m *mockDataStore) WasSetFailedCalled() bool {
	m.setFailedMutex.Lock()
	defer m.setFailedMutex.Unlock()
	return m.wasSetFailedCalled
}

func (m *mockDataStore) WasSetFinishedCalled() bool {
	m.setFinishedMutex.Lock()
	defer m.setFinishedMutex.Unlock()
	return m.wasSetFinishedCalled
}

func (m *mockDataStore) WasStoreCalled() bool {
	m.storeMutex.Lock()
	defer m.storeMutex.Unlock()
	return m.wasStoreCalled
}

type mockQueue struct {
	failJob bool

	finishChan chan []string
	failChan   chan []string
}

func (m *mockQueue) Enqueue(url string, destination string) error {
	if m.failJob {
		m.failChan <- []string{url, destination}
	} else {
		m.finishChan <- []string{url, destination}
	}
	// to simulate job proccessing
	time.Sleep(1 * time.Second)
	return nil
}

func (m *mockQueue) FinishChan() <-chan []string {
	return m.finishChan
}

func (m *mockQueue) FailChan() <-chan []string {
	return m.failChan
}
