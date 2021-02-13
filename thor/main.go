package main

import (
	"context"
	"log"
	"net"
	"os"

	"cloud.google.com/go/storage"
	pb "github.com/2hamed/saas/protobuf"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

type server struct {
	pb.CaptureServer

	capture Capture
}

func (s *server) Capture(ctx context.Context, in *pb.CaptureRequest) (*pb.CaptureResponse, error) {

	return nil, nil
}

func main() {
	port := os.Getenv("GRPC_LISTEN_PORT")

	serviceAccountFile := os.Getenv("GCP_SERVICE_ACCOUNT_FILE_PATH")
	stoageClient, err := storage.NewClient(context.Background(), option.WithServiceAccountFile(serviceAccountFile))
	if err != nil {
		panic(err)
	}

	bucketName := os.Getenv("GCP_STORAGE_BUCKET_NAME")
	capture, err := NewCapture(stoageClient, bucketName)
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCaptureServer(s, &server{
		capture: capture,
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
