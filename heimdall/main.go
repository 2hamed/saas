package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	pb "github.com/2hamed/saas/protobuf"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

func main() {
	listenHostPort := fmt.Sprintf(":%s", os.Getenv("HTTP_LISTEN_PORT"))
	address := os.Getenv("QUEUE_GRPC_ADDRESS")

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	queueClient := pb.NewQueueClient(conn)

	h := &Handler{queueClient}

	r := chi.NewRouter()

	r.Post("/", h.NewJob)

	log.Println("HTTP server listening on", listenHostPort)

	if err := http.ListenAndServe(listenHostPort, r); err != nil {
		log.Fatalf("HTTP server failed to listen: %v", err)
	}
}
