package filestore

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/2hamed/saas/waitfor"
	"github.com/minio/minio-go/v6"
	log "github.com/sirupsen/logrus"
)

const (
	bucketName = "shots"
	location   = "global"
)

var (
	minioHost string
	minioPort string
)

func initConfig() {
	minioHost = os.Getenv("MINIO_HOST")
	minioPort = os.Getenv("MINIO_PORT")
}

func createMinio() (*minioFileStore, error) {

	initConfig()

	waitfor.WaitForServices([]string{
		fmt.Sprintf("%s:%s", minioHost, minioPort),
	}, 60*time.Second)

	log.Infof("Connecting to Minio on %s:%s", minioHost, minioPort)

	endpoint := fmt.Sprintf("%s:%s", minioHost, minioPort)
	accessKeyID := os.Getenv("MINIO_ACCESSKEY")
	secretAccessKey := os.Getenv("MINIO_SECRET")

	useSSL := false

	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, err
	}

	err = client.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := client.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			log.Debugf("We already own %s", bucketName)
		} else {
			return nil, fmt.Errorf("failed to create minio bucket: %w, %v", err, errBucketExists)
		}
	} else {
		log.Debugf("Successfully created bucket[%s]", bucketName)
	}

	policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":["arn:aws:s3:::shots"]},{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::shots/*"]}]}`

	err = client.SetBucketPolicy(bucketName, policy)
	if err != nil {
		return nil, err
	}

	log.Debugf("Minio bucket policy: %s", policy)

	return &minioFileStore{
		client: client,
	}, nil
}

type minioFileStore struct {
	client *minio.Client
}

func (m *minioFileStore) Store(path string) (string, error) {
	log.Debug("Uploading file to minio")

	contentType, err := getFileContentType(path)

	if err != nil {
		return "", fmt.Errorf("failed get file content type: %w", err)
	}

	fileName := getFileName(path)

	_, err = m.client.FPutObject(bucketName, fileName, path, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", bucketName, fileName), nil
}

func getFileContentType(filePath string) (string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, 512)

	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func getFileName(path string) string {
	parts := strings.Split(path, "/")

	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}
