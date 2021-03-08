package trace

import (
	"context"
	"fmt"
	"os"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	defaultTracer trace.Tracer
	projectID     string

	DefaultTracerName = "saas"
)

func StartTracing(opts ...TraceOpts) error {
	tracingConfig := &tracingConfig{}
	for _, opt := range opts {
		opt(tracingConfig)
	}

	DefaultTracerName = tracingConfig.traceName

	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	exporter, err := texporter.NewExporter(texporter.WithProjectID(projectID))
	if err != nil {
		return err
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	otel.SetTracerProvider(tp)
	return nil
}
func GetTracer() trace.Tracer {
	if defaultTracer == nil {
		defaultTracer = otel.Tracer(DefaultTracerName)
	}
	return defaultTracer
}

func Start(ctx context.Context, name string, opts ...trace.SpanOption) (context.Context, trace.Span) {
	return GetTracer().Start(ctx, name)
}

func FQDN(span trace.Span) string {
	return fmt.Sprintf("projects/%s/traces/%s", projectID, span.SpanContext().TraceID.String())
}
