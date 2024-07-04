package gcplogurl

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
	octrace "go.opencensus.io/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TraceLogURL construct Cloud Logging URL about this trace (request).
func TraceLogURL(ctx context.Context) (string, error) {
	projectID, _ := metadata.ProjectID()
	if projectID == "" {
		for _, key := range []string{"GCP_PROJECT", "GCLOUD_PROJECT", "GOOGLE_CLOUD_PROJECT"} {
			projectID = os.Getenv(key)
			if projectID != "" {
				break
			}
		}
	}
	if projectID == "" {
		return "", errors.New("failed to detect GCP Project ID")
	}

	var traceID string

	traceLib := os.Getenv("GOOGLE_API_GO_EXPERIMENTAL_TELEMETRY_PLATFORM_TRACING")
	switch traceLib {
	case "", "opencensus":
		span := octrace.FromContext(ctx)
		if span == nil {
			return "", errors.New("ctx doesn't have OpenCensus span")
		}
		traceID = span.SpanContext().TraceID.String()

	case "opentelemetry":
		span := oteltrace.SpanFromContext(ctx)
		if span == nil {
			return "", errors.New("ctx doesn't have OpenTelemetry span")
		}
		traceID = span.SpanContext().TraceID().String()

	default:
		return "", fmt.Errorf("unknown trace lib %q", traceLib)
	}

	ex := &Explorer{
		ProjectID: projectID,
		Query:     Query(fmt.Sprintf(`trace="projects/%s/traces/%s"`, projectID, traceID)),
		TimeRange: &SpecificTimeWithRange{
			At:    time.Now(),
			Range: 2 * time.Hour,
		},
	}
	return ex.String(), nil
}
