package dispatcher

import (
	"crypto/sha1"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// Dispatcher interface is responsible for dispatching jobs to job queue and returning the result to caller
type Dispatcher interface {
	Enqueue(url string) error

	FetchResult(url string) (string, error)

	FetchStatus(url string) (exists bool, isFinished bool, err error)
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

	fileName := hash(url) + ".png"

	fullPath := fmt.Sprintf("%s/%s", os.Getenv("STORAGE_PATH"), fileName)

	err := d.ds.Store(url, fileName)

	if err != nil {
		return fmt.Errorf("failed to store to DataStore: %v", err)
	}

	return d.q.Enqueue(url, fullPath)
}

func (d *dispatcher) FetchStatus(url string) (exists bool, isFinished bool, err error) {
	return d.ds.FetchStatus(url)
}

func (d *dispatcher) FetchResult(url string) (string, error) {
	return d.ds.Fetch(url)
}

func (d *dispatcher) listenForResults() {

	go func() {
		for {
			select {
			case r := <-d.q.FinishChan():
				log.Debug("Received finish signal for", r)

				err := d.ds.SetFinished(r[0])

				if err != nil {
					log.Error("failed to update status for", r, err)
				}
			case r := <-d.q.FailChan():
				log.Debug("Received fail signal for ", r)

				err := d.ds.SetFailed(r[0])

				if err != nil {
					log.Error("failed to set failed for", r, err)
				}
			}
		}
	}()
}

func hash(url string) string {
	hash := sha1.New()
	_, _ = hash.Write([]byte(url))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
