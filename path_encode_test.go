package gcplogurl

import (
	"testing"
)

func Test_values(t *testing.T) {
	tests := []struct {
		name          string
		v             values
		wantEncode    string
		wantEncodeRaw string
	}{
		{
			name: "basic",
			v: values{
				"storageScope": []string{"project"},
			},
			wantEncode:    "storageScope=project",
			wantEncodeRaw: "storageScope=project",
		},
		{
			name: "basic - query",
			v: values{
				"query": []string{`trace="projects/test-project/traces/f7de19bd944a9229a03f165a4f8cb092"`},
			},
			wantEncode:    "query=trace%3D%22projects%2Ftest-project%2Ftraces%2Ff7de19bd944a9229a03f165a4f8cb092%22",
			wantEncodeRaw: `query=trace="projects/test-project/traces/f7de19bd944a9229a03f165a4f8cb092"`,
		},
		{
			name: "multi values",
			v: values{
				"summaryFields": []string{"trace", "protoPayload/resource"},
			},
			wantEncode:    "summaryFields=trace,protoPayload%2Fresource",
			wantEncodeRaw: "summaryFields=trace,protoPayload/resource",
		},
		{
			name: "multi keys",
			v: values{
				"summaryFields": []string{"trace", "protoPayload/resource"},
				"storageScope":  []string{"storage", "projects/test-project/locations/global/buckets/_Default/views/_AllLogs"},
			},
			wantEncode:    "storageScope=storage,projects%2Ftest-project%2Flocations%2Fglobal%2Fbuckets%2F_Default%2Fviews%2F_AllLogs;summaryFields=trace,protoPayload%2Fresource",
			wantEncodeRaw: "storageScope=storage,projects/test-project/locations/global/buckets/_Default/views/_AllLogs;summaryFields=trace,protoPayload/resource",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Encode(); got != tt.wantEncode {
				t.Errorf("Encode() = %v, want %v", got, tt.wantEncode)
			}
			if got := tt.v.RawEncode(); got != tt.wantEncodeRaw {
				t.Errorf("RawEncode() = %v, want %v", got, tt.wantEncodeRaw)
			}
		})
	}
}
