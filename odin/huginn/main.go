package main

import (
	"context"
	"log"
	"net"

	odin "github.com/2hamed/saas/odin"
	pb "github.com/2hamed/saas/protobuf"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

const port = ":50005"

type server struct {
	pb.QueueServer

	q odin.QManager
}

func (s *server) Capture(ctx context.Context, in *pb.QueueRequest) (*pb.QueueResponse, error) {
	uuid := uuid.New()
	err := s.q.Enqueue(ctx, odin.CaptureJob{
		UUID: uuid.String(),
		URL:  in.GetUrl(),
	})
	return &pb.QueueResponse{Uuid: uuid.String()}, err
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	qManager, err := odin.NewQManager()
	if err != nil {
		log.Fatalf("failed creating queue manger: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterQueueServer(s, &server{q: qManager})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
