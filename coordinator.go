package main

import (
	"github.com/2hamed/saas/jobq"
	"github.com/2hamed/saas/storage"
)

// Coordinator is the highlevel component gluing all the other components together
type Coordinator struct {
	jobQ        jobq.QManager
	dataStore   storage.DataStore
	dataFetcher storage.DataFetcher
}

// CaptureAsync asynchronously captures the screenshot of the specified url
func (c *Coordinator) CaptureAsync(url string) error {
	err := c.dataStore.Store("", "")
	if err != nil {
		return err
	}

	err = c.jobQ.Enqueue(url)

	return err
}

func (c *Coordinator) JobFailed(url string) {

}

func (c *Coordinator) JobFinished(url string) {

}
