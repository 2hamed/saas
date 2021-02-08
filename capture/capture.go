package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	log "github.com/sirupsen/logrus"
)

// Capture is the interface to capture the screenshot
type Capture interface {
	Save(ctx context.Context, url string) (string, error)
}

// NewCapture returns a concrete implementaion of the Capture interface
func NewCapture(storageClient *storage.Client, bucketName string) (Capture, error) {
	return phantomJs{storageClient, bucketName}, nil
}

type phantomJs struct {
	storageClient *storage.Client
	bucketName    string
}

func (p phantomJs) Save(ctx context.Context, url string) (string, error) {
	log.Debugf("Capture initiated for %s", url)

	hash := md5.Sum([]byte(url))
	path := "/tmp/" + string(hash[:])

	cmd := exec.Command(os.Getenv("PHANTOMJS_PATH"), os.Getenv("CAPTUREJS_PATH"), url, path)

	log.Debugf("Executing %s", cmd.String())

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("error capturing screenshot: %s (%w)", output, err)
	}

	objectName := string(hash[:]) + ".png"
	return objectName, p.uploadToStorage(ctx, path, objectName)
}
func (p phantomJs) uploadToStorage(ctx context.Context, path, name string) error {
	w := p.storageClient.Bucket(p.bucketName).Object(name).NewWriter(ctx)

	f, err := os.OpenFile(path, os.O_RDONLY, os.FileMode(755))
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	return err

}
