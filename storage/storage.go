package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/2hamed/saas/waitfor"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DataStore is the abstract type to work with the underlying datastore
type DataStore interface {
	Fetch(url string) (string, error)
	FetchStatus(url string) (exists bool, isPending bool, isFinished bool, err error)

	Store(url string, destination string) error
	SetFinished(url string) error
	SetFailed(url string) error
}

// NewDataStore returns an implementation of DataStore
func NewDataStore() (DataStore, error) {

	waitfor.WaitForServices([]string{
		fmt.Sprintf("%s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT")),
	}, 10*time.Second)

	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT")))

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to Mongo: %v", err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 1*time.Second)
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed pinging Mongo: %v", err)
	}

	log.Infof("Conncted to MongoDB on %s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"))

	return &mongoDataStore{
		client: client,
	}, nil
}
