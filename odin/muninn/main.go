package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	odin "github.com/2hamed/saas/odin"
	pb "github.com/2hamed/saas/protobuf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type muninn struct {
	q          odin.QManager
	grpcClient pb.CaptureClient
}

func main() {
	log.Logger = log.Level(zerolog.InfoLevel)
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log.Info().Msg("Starting Muninn...")

	address := os.Getenv("CAPTURE_GRPC_ADDRESS")

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal().Err(err).Msg("did not connect")
	}
	defer conn.Close()

	grpcClient := pb.NewCaptureClient(conn)
	qManager, err := odin.NewQManager()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create q manager")
	}

	log.Info().Msg("Muninn is waiting for jobs...")
	muninn := &muninn{
		q:          qManager,
		grpcClient: grpcClient,
	}
	muninn.Start()

	log.Info().Msg("Shutting down...")

}

func (m *muninn) Start() error {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobChan, err := m.q.GetJobChan(ctx)
	if err != nil {
		return err
	}

loop:
	for {
		select {
		case <-sigChan:
			break loop
		case job := <-jobChan:
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			log.Info().Str("url", job.URL).Str("uuid", job.UUID).Msg("Received job")
			resp, err := m.grpcClient.Capture(ctx, &pb.CaptureRequest{
				Uuid: job.UUID,
				Url:  job.URL,
			})
			if err != nil {
				log.Error().Err(err).Msg("Error making capture")
				job.Nack()
			} else {
				log.Info().Str("object_path", resp.ObjectPath).Str("uuid", resp.Uuid).Msg("Screenshot successfully captured!")
				job.Ack()
			}
			cancel()
		}
	}

	m.q.CleanUp()
	return nil
}
