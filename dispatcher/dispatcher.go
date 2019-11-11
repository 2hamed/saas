package dispatcher

import (
	"crypto/sha1"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Dispatcher interface is responsible for dispatching jobs to job queue and returning the result to caller
type Dispatcher interface {
	Enqueue(url string) error

	GetResult(url string) (string, error)
}

// NewDispatcher returns a new instance of Dispatcher interface
func NewDispatcher(ds dataStore, q queue) (Dispatcher, error) {
	d := &dispatcher{
		ds: ds,
		q:  q,
	}

	d.listenForResults()

	return d, nil
}

type dispatcher struct {
	ds dataStore
	q  queue
}

func (d *dispatcher) Enqueue(url string) error {

	log.Debug("Received url to process", url)

	destPath := getSha1(url)

	err := d.ds.Store(url, destPath)

	if err != nil {
		return fmt.Errorf("failed to store to DataStore: %v", err)
	}

	return d.q.Enqueue(url, fmt.Sprintf("/tmp/%s.png", destPath))
}

func (d *dispatcher) GetResult(url string) (string, error) {

	return "", nil
}

func (d *dispatcher) listenForResults() {
	go func() {
		for {
			result := <-d.q.FinishChan()

			log.Debug("Received finish signal for", result)

			err := d.ds.UpdateStatus(result[0], true)

			if err != nil {
				log.Error("failed to update status for", result, err)
			}
		}
	}()

	go func() {
		for {
			result := <-d.q.FailChan()

			log.Debug("Received fail signal for ", result)

			err := d.ds.SetFailed(result[0])

			if err != nil {
				log.Error("failed to set failed for", result, err)
			}

		}
	}()
}

func getSha1(url string) string {
	hash := sha1.New()
	hash.Write([]byte(url))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
