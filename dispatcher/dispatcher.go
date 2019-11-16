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

	FetchStatus(url string) (exists bool, isPending bool, isFinished bool, err error)
}

// NewDispatcher returns a new instance of Dispatcher interface
// ds (dataStore) and q (queue) are required
// fs (fileStore) is optional and can be nil
func NewDispatcher(ds dataStore, q queue, fs filestore) (Dispatcher, error) {
	d := &dispatcher{
		ds: ds,
		q:  q,
		fs: fs,
	}

	d.listenForResults()

	return d, nil
}

type dispatcher struct {
	ds dataStore
	q  queue
	fs filestore
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

func (d *dispatcher) FetchStatus(url string) (exists bool, isPending bool, isFinished bool, err error) {
	return d.ds.FetchStatus(url)
}

func (d *dispatcher) FetchResult(url string) (string, error) {
	return d.ds.Fetch(url)
}

func (d *dispatcher) listenForResults() {
	go func() {
		for {
			select {
			case result := <-d.q.FinishChan():

				log.Debugf("Received finish signal for: %v", result)

				d.processFinishSignal(result)

			case result := <-d.q.FailChan():

				log.Debugf("Received fail signal for: %v", result)

				d.processFailSignal(result)
			}
		}
	}()
}

func (d *dispatcher) processFinishSignal(result []string) {
	if d.fs != nil {
		remoteURL, err := d.fs.Store(result[1])

		if err != nil {
			log.Errorf("failed to store file in filestore: %v", err)
		} else {
			err = d.ds.UpdatePath(result[0], remoteURL)
			if err != nil {
				log.Errorf("failed to update path for url: %v", err)
			}
		}
	}

	err := d.ds.SetFinished(result[0])

	if err != nil {
		log.Error("failed to update status for ", result, err)
	}
}

func (d *dispatcher) processFailSignal(result []string) {
	err := d.ds.SetFailed(result[0])

	if err != nil {
		log.Error("failed to set failed for", result, err)
	}
}

func hash(url string) string {
	hash := sha1.New()
	_, _ = hash.Write([]byte(url))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
