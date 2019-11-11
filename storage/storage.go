package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataStore interface {
	Fetch(url string) (string, error)
	Store(url string, destination string) error
	UpdateStatus(url string, finished bool) error
	SetFailed(url string) error
}

var (
	mongoHost = os.Getenv("MONGO_HOST")
	mongoPort = os.Getenv("MONGO_PORT")
)

func NewDataStore() (DataStore, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)))
	if err != nil {
		return nil, fmt.Errorf("failed connecting to Mongo: %v", err)
	}
	return &mongoDataStore{
		client: client,
	}, nil
}

const (
	databaseName   = "screenshots"
	collectionName = "jobs"
)

type mongoDataStore struct {
	client *mongo.Client
}

func (s *mongoDataStore) Fetch(url string) (string, error) {
	res := s.client.Database(databaseName).Collection(collectionName).FindOne(context.Background(), bson.M{
		"url": url,
	})

	var m bson.M
	err := res.Decode(&m)
	if err != nil {
		return "", err
	}
	return m["path"].(string), nil
}

func (s *mongoDataStore) Store(url string, destination string) error {
	_, err := s.client.Database(databaseName).Collection(collectionName).InsertOne(context.Background(), bson.M{
		"url":    url,
		"path":   destination,
		"status": "pending",
		"time":   time.Now().Unix(),
	})

	return err
}

func (s *mongoDataStore) UpdateStatus(url string, finished bool) error {
	status := "finished"
	if !finished {
		status = "pending"
	}
	_, err := s.client.Database(databaseName).Collection(collectionName).UpdateOne(context.Background(),
		bson.M{
			"url": url,
		},
		bson.M{
			"$set": bson.M{
				"status": status,
			},
		})

	return err
}

func (s *mongoDataStore) SetFailed(url string) error {
	_, err := s.client.Database(databaseName).Collection(collectionName).UpdateOne(context.Background(),
		bson.M{
			"url": url,
		},
		bson.M{
			"status": "failed",
		})
	return err
}
