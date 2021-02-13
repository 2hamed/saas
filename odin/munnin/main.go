package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	odin "github.com/2hamed/saas/odin"
	pb "github.com/2hamed/saas/protobuf"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

type muninn struct {
	q          odin.QManager
	grpcClient pb.CaptureClient
}

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	grpcClient := pb.NewCaptureClient(conn)
	qManager, err := odin.NewQManager()
	if err != nil {
		log.Fatalf("failed to create q manager: %v", err)
	}

	log.Println("Muninn is waiting for jobs...")
	muninn := &muninn{
		q:          qManager,
		grpcClient: grpcClient,
	}
	muninn.Start()

	log.Println("Shutting down...")

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
			m.grpcClient.Capture(ctx, &pb.CaptureRequest{
				Uuid: job.UUID,
				Url:  job.URL,
			})
		}
	}

	m.q.CleanUp()
	return nil
}
