package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	odin "github.com/2hamed/saas/odin"
	pb "github.com/2hamed/saas/protobuf"
	"github.com/2hamed/saas/trace"
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

	if err := trace.StartTracing(trace.WithName("muninn")); err != nil {
		log.Fatal().Err(err).Msg("failed to intialize tracing")
	}

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
	defer m.q.CleanUp()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobChan, err := m.q.GetJobChan(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-sigChan:
			return nil
		case job := <-jobChan:
			m.processJob(ctx, job)
		}
	}

}
func (m *muninn) processJob(ctx context.Context, job odin.CaptureJob) {
	ctx, span := trace.Start(ctx, "Capture")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	log.Info().Str("trace", trace.FQDN(span)).Str("span_id", span.SpanContext().SpanID.String()).Str("url", job.URL).Str("uuid", job.UUID).Msg("Received job")
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
}
