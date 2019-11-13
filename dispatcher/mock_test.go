package dispatcher

type mockDataStore struct {
	fetchResult string
	fetchErr    error

	fetchStatusExists     bool
	fetchStatusIsFinished bool
	fetchStatusErr        error

	storeErr error

	setFinishedErr error
	setFailedErr   error

	wasFetchCalled       bool
	wasFetchStatusCalled bool
	wasStoreCalled       bool
	wasSetFailedCalled   bool
	wasSetFinishedCalled bool
}

func (m *mockDataStore) Fetch(url string) (string, error) {
	m.wasFetchCalled = true
	return m.fetchResult, m.fetchErr
}

func (m *mockDataStore) FetchStatus(url string) (exists bool, isFinished bool, err error) {
	m.wasFetchStatusCalled = true
	return m.fetchStatusExists, m.fetchStatusIsFinished, m.fetchStatusErr
}

func (m *mockDataStore) Store(url string, destination string) error {
	m.wasStoreCalled = true
	return m.storeErr
}

func (m *mockDataStore) SetFinished(url string) error {
	m.wasSetFinishedCalled = true
	return m.setFinishedErr
}

func (m *mockDataStore) SetFailed(url string) error {
	m.wasSetFailedCalled = true
	return m.setFailedErr
}

type mockQueue struct {
	failJob bool

	finishChan chan []string
	failChan   chan []string
}

func (m *mockQueue) Enqueue(url string, destination string) error {
	if m.failJob {
		m.failChan <- []string{}
	} else {
		m.finishChan <- []string{}
	}
	return nil
}

func (m *mockQueue) FinishChan() <-chan []string {
	return m.finishChan
}

func (m *mockQueue) FailChan() <-chan []string {
	return m.failChan
}
