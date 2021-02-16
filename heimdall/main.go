package main

import (
	"fmt"
	"net/http"
	"os"

	pb "github.com/2hamed/saas/protobuf"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	log.Logger = log.Level(zerolog.InfoLevel)

	listenHostPort := fmt.Sprintf(":%s", os.Getenv("HTTP_LISTEN_PORT"))
	address := os.Getenv("QUEUE_GRPC_ADDRESS")

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal().Err(err).Msg("grpc client did not connect")
	}
	defer conn.Close()

	queueClient := pb.NewQueueClient(conn)

	h := &Handler{queueClient}

	r := chi.NewRouter()

	r.Post("/", h.NewJob)

	log.Info().Msgf("HTTP server listening on: %s", listenHostPort)

	if err := http.ListenAndServe(listenHostPort, r); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed to listen")
	}
}
