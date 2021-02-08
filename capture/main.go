package main

import (
	"context"
	"log"
	"net"

	pb "github.com/2hamed/saas/protobuf"
	"google.golang.org/grpc"
)

const port = ":50005"

type server struct {
	pb.CaptureServer
}

func (s *server) Capture(ctx context.Context, in *pb.CaptureRequest) (*pb.CaptureResponse, error) {
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCaptureServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
