package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/rs/zerolog/log"
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
	log.Info().Str("url", url).Msg("Capture initiated for")

	hash := md5.Sum([]byte(url))
	path := fmt.Sprintf("/tmp/%x.png", hash)

	cmd := exec.Command(os.Getenv("PHANTOMJS_PATH"), os.Getenv("CAPTUREJS_PATH"), url, path)

	log.Debug().Str("cmd", cmd.String()).Msg("Capture excuting")

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("error capturing screenshot: %s (%w)", output, err)
	}

	objectName := fmt.Sprintf("%x.png", hash)
	return objectName, p.uploadToStorage(ctx, path, objectName)
}
func (p phantomJs) uploadToStorage(ctx context.Context, path, name string) error {
	f, err := os.OpenFile(path, os.O_RDONLY, os.FileMode(755))
	if err != nil {
		return err
	}
	defer f.Close()

	w := p.storageClient.Bucket(p.bucketName).Object(name).NewWriter(ctx)
	defer w.Close()

	_, err = io.Copy(w, f)
	return err

}
