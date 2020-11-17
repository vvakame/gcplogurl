package gcplogurl_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opencensus.io/trace"

	"github.com/vvakame/gcplogurl"
)

func Example() {
	exp := &gcplogurl.Explorer{
		ProjectID:    "test-project",
		StorageScope: gcplogurl.StorageScopeProject,
		Query:        `trace="projects/test-project/traces/2f029f622a49837d621cbe6dc5e5d146"`,
		TimeRange: &gcplogurl.SpecificTimeWithRange{
			At:    time.Date(2020, 11, 14, 4, 20, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
			Range: 2 * time.Hour,
		},
	}
	logURL := exp.String()

	fmt.Println(logURL)

	// Output:
	// https://console.cloud.google.com/logs/query;query=trace%3D%22projects%2Ftest-project%2Ftraces%2F2f029f622a49837d621cbe6dc5e5d146%22;storageScope=project;timeRange=2020-11-13T19:20:00Z%2F2020-11-13T19:20:00Z--PT2H?project=test-project
}

func ExampleTraceLogURL() {
	_ = os.Setenv("GCP_PROJECT", "test-project")
	ctx := context.Background()
	ctx, span := trace.StartSpan(ctx, "test-span")
	defer span.End()

	logURL, err := gcplogurl.TraceLogURL(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(logURL)
}
