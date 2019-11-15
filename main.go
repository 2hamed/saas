package main

import (
	"github.com/2hamed/saas/api"
	"github.com/2hamed/saas/dispatcher"
	"github.com/2hamed/saas/filestore"
	"github.com/2hamed/saas/jobq"
	"github.com/2hamed/saas/screenshot"
	"github.com/2hamed/saas/storage"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {

	err := godotenv.Load("/home/hamed/dev/projects/saas/.env")
	if err != nil {
		log.Warn("Error loading .env file")
	}

	log.SetLevel(log.DebugLevel)

	capture, err := screenshot.NewCapture()
	if err != nil {
		log.Fatal(err)
	}

	qManager, err := jobq.NewQManager(capture)
	if err != nil {
		log.Fatal(err)
	}

	ds, err := storage.NewDataStore()
	if err != nil {
		log.Fatal(err)
	}

	fs, err := filestore.NewFileStore()
	if err != nil {
		log.Fatal(err)
	}

	dispatcher, err := dispatcher.NewDispatcher(ds, qManager, fs)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(api.StartServer(dispatcher))

}
