package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"cloud.google.com/go/storage"
	pb "github.com/2hamed/saas/protobuf"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

type server struct {
	pb.CaptureServer

	capture Capture
}

func (s *server) Capture(ctx context.Context, in *pb.CaptureRequest) (*pb.CaptureResponse, error) {
	objectPath, err := s.capture.Save(ctx, in.Url)
	if err != nil {
		return nil, err
	}

	return &pb.CaptureResponse{
		Uuid:       in.Uuid,
		ObjectPath: objectPath,
	}, nil
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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	s := grpc.NewServer()
	pb.RegisterCaptureServer(s, &server{
		capture: capture,
	})

	log.Info().Msgf("GRPC Server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
