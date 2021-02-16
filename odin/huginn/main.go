package main

import (
	"context"
	"fmt"
	"net"
	"os"

	odin "github.com/2hamed/saas/odin"
	pb "github.com/2hamed/saas/protobuf"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

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

	log.Logger = log.Level(zerolog.InfoLevel)
	port := os.Getenv("GRPC_LISTEN_PORT")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	qManager, err := odin.NewQManager()
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating queue manger")
	}

	s := grpc.NewServer()
	pb.RegisterQueueServer(s, &server{q: qManager})

	log.Info().Msgf("GRPC Server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
