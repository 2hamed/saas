package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	pb "github.com/2hamed/saas/protobuf"
	"github.com/2hamed/saas/trace"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	log.Logger = log.Level(zerolog.InfoLevel)
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log.Info().Msg("Starting Heimdall...")

	if err := trace.StartTracing(trace.WithName("heimdall")); err != nil {
		log.Fatal().Err(err).Msg("failed to intialize tracing")
	}

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

	r.Get("/", serveHome)
	r.Post("/c", h.NewJob)

	log.Info().Msgf("HTTP server listening on: %s", listenHostPort)

	if err := http.ListenAndServe(listenHostPort, r); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed to listen")
	}
}
func serveHome(w http.ResponseWriter, r *http.Request) {
	_, span := trace.Start(r.Context(), "Home")
	defer span.End()

	w.WriteHeader(200)
	w.Write([]byte("Hello! I am Heimdall!"))
}
