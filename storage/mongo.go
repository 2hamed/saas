package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/2hamed/saas/waitfor"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "screenshots"
	collectionName = "jobs"
)

var (
	mongoHost string
	mongoPort string
)

func initConfig() {
	mongoHost = os.Getenv("MONGO_HOST")
	mongoPort = os.Getenv("MONGO_PORT")
}

type mongoDataStore struct {
	client *mongo.Client
}

func createMongoDataStore() (*mongoDataStore, error) {

	initConfig()

	waitfor.WaitForServices([]string{
		fmt.Sprintf("%s:%s", mongoHost, mongoPort),
	}, 10*time.Second)

	log.Infof("Conncting to MongoDB on %s:%s", mongoHost, mongoPort)

	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort))

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to Mongo: %w", err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 1*time.Second)
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed pinging Mongo: %w", err)
	}

	return &mongoDataStore{
		client: client,
	}, nil
}

func (s *mongoDataStore) Fetch(url string) (string, error) {
	res := s.client.Database(databaseName).Collection(collectionName).FindOne(context.Background(), bson.M{
		"url": url,
	})
	if res.Err() != nil {
		return "", fmt.Errorf("failed fetching path from database: %w", res.Err())
	}
	var m bson.M
	err := res.Decode(&m)
	if err != nil {
		return "", fmt.Errorf("failed decoding object from database: %w", err)
	}
	return m["path"].(string), nil
}

func (s *mongoDataStore) FetchStatus(url string) (exists bool, isPending bool, isFinished bool, err error) {
	res := s.client.Database(databaseName).Collection(collectionName).FindOne(context.Background(), bson.M{
		"url": url,
	})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return false, false, false, nil
		}
		return false, false, false, fmt.Errorf("failed fetching path from database: %w", res.Err())
	}

	var m bson.M
	err = res.Decode(&m)
	if err != nil {
		return false, false, false, fmt.Errorf("failed decoding object from database: %w", err)
	}
	if m["status"].(string) == "pending" {
		return true, true, false, nil
	} else if m["status"].(string) == "finished" {
		return true, false, true, nil
	} else if m["status"].(string) == "failed" {
		return true, false, false, nil
	}

	return false, false, false, errors.New("this is an unkown status and should never happen")
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

func (s *mongoDataStore) SetFinished(url string) error {

	_, err := s.client.Database(databaseName).Collection(collectionName).UpdateOne(context.Background(),
		bson.M{
			"url": url,
		},
		bson.M{
			"$set": bson.M{
				"status": "finished",
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
			"$set": bson.M{
				"status": "failed",
			},
		})
	return err
}

func (s *mongoDataStore) Delete(url string) error {
	_, err := s.client.Database(databaseName).Collection(collectionName).DeleteMany(context.Background(), bson.M{
		"url": url,
	})
	return err
}

func (s *mongoDataStore) UpdatePath(url string, destination string) error {
	_, err := s.client.Database(databaseName).Collection(collectionName).UpdateOne(context.Background(),
		bson.M{
			"url": url,
		},
		bson.M{
			"$set": bson.M{
				"path": destination,
			},
		})
	return err
}
