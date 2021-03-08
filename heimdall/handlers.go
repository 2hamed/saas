package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	pb "github.com/2hamed/saas/protobuf"
	"github.com/2hamed/saas/trace"
	"github.com/rs/zerolog/log"
)

type CaptureRequest struct {
	URL string `json:"url"`
}

type Handler struct {
	grpcClient pb.QueueClient
}

func (h *Handler) NewJob(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.Start(r.Context(), "Capture")
	defer span.End()

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Error().Str("trace", trace.FQDN(span)).Str("span_id", span.SpanContext().SpanID.String()).Err(err).Msg("Error reading request body")
		httpErr(w, 500, "error occurred")
		return
	}
	var captureRequest CaptureRequest
	err = json.Unmarshal(body, &captureRequest)
	if err != nil {
		log.Error().Str("trace", trace.FQDN(span)).Str("span_id", span.SpanContext().SpanID.String()).Err(err).Msg("Error unmarshaling request body")
		httpErr(w, 500, "error occurred")
		return
	}

	log.Info().Str("trace", trace.FQDN(span)).Str("span_id", span.SpanContext().SpanID.String()).Str("url", captureRequest.URL).Msg("Received request to capture")

	resp, err := h.grpcClient.Capture(ctx, &pb.QueueRequest{
		Url: captureRequest.URL,
	})
	if err != nil {
		log.Error().Str("trace", trace.FQDN(span)).Str("span_id", span.SpanContext().SpanID.String()).Err(err).Msg("Error making grpc request")
		httpErr(w, 500, "error occurred")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("{\"uuid\":\"%s\"}", resp.GetUuid())))

}

func httpErr(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}
