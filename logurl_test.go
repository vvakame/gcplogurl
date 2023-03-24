package gcplogurl_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/vvakame/gcplogurl"
)

func TestExplorer(t *testing.T) {
	tests := []struct {
		name    string
		obj     *gcplogurl.Explorer
		want    string
		wantErr bool
	}{
		// Query
		{
			name: "query",
			obj: &gcplogurl.Explorer{
				ProjectID: "test-project",
				Query:     `trace="projects/test-project/traces/2f029f622a49837d621cbe6dc5e5d146"`,
			},
			want: "https://console.cloud.google.com/logs/query;query=trace%3D%22projects%2Ftest-project%2Ftraces%2F2f029f622a49837d621cbe6dc5e5d146%22?project=test-project",
		},
		// StorageScope
		{
			name: "storage scope - project",
			obj: &gcplogurl.Explorer{
				ProjectID:    "test-project",
				StorageScope: gcplogurl.StorageScopeProject,
			},
			want: "https://console.cloud.google.com/logs/query;storageScope=project?project=test-project",
		},
		{
			name: "storage scope - storage",
			obj: &gcplogurl.Explorer{
				ProjectID: "test-project",
				StorageScope: &gcplogurl.StorageScopeStorage{
					Src: []string{
						"projects/test-project/locations/global/buckets/_Default/views/_AllLogs",
						"projects/test-project/locations/global/buckets/_Default/views/_Default",
					},
				},
			},
			want: "https://console.cloud.google.com/logs/query;storageScope=storage,projects%2Ftest-project%2Flocations%2Fglobal%2Fbuckets%2F_Default%2Fviews%2F_AllLogs,projects%2Ftest-project%2Flocations%2Fglobal%2Fbuckets%2F_Default%2Fviews%2F_Default?project=test-project",
		},
		// TimeRange
		{
			name: "time range - last 30 min ago",
			obj: &gcplogurl.Explorer{
				ProjectID: "test-project",
				TimeRange: &gcplogurl.RecentRange{
					Last: 30 * time.Minute,
				},
			},
			want: "https://console.cloud.google.com/logs/query;timeRange=PT30M?project=test-project",
		},
		{
			name: "time range - specific time between",
			obj: &gcplogurl.Explorer{
				ProjectID: "test-project",
				TimeRange: &gcplogurl.SpecificTimeBetween{
					From: time.Date(2020, 11, 13, 4, 20, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
					To:   time.Date(2020, 11, 14, 4, 20, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				},
			},
			want: "https://console.cloud.google.com/logs/query;timeRange=2020-11-12T19:20:00Z%2F2020-11-13T19:20:00Z?project=test-project",
		},
		{
			name: "time range - specific time from (only)",
			obj: &gcplogurl.Explorer{
				ProjectID: "test-project",
				TimeRange: &gcplogurl.SpecificTimeBetween{
					From: time.Date(2020, 11, 13, 4, 20, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				},
			},
			want: "https://console.cloud.google.com/logs/query;timeRange=2020-11-12T19:20:00Z%2F?project=test-project",
		},
		{
			name: "time range - specific time between to (only)",
			obj: &gcplogurl.Explorer{
				ProjectID: "test-project",
				TimeRange: &gcplogurl.SpecificTimeBetween{
					To: time.Date(2020, 11, 14, 4, 20, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				},
			},
			want: "https://console.cloud.google.com/logs/query;timeRange=%2F2020-11-13T19:20:00Z?project=test-project",
		},
		{
			name: "time range - specific time with range",
			obj: &gcplogurl.Explorer{
				ProjectID: "test-project",
				TimeRange: &gcplogurl.SpecificTimeWithRange{
					At:    time.Date(2020, 11, 14, 4, 20, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
					Range: 2 * time.Hour,
				},
			},
			want: "https://console.cloud.google.com/logs/query;timeRange=2020-11-13T19:20:00Z%2F2020-11-13T19:20:00Z--PT2H?project=test-project",
		},
		// SummaryFields
		{
			name: "summary fields",
			obj: &gcplogurl.Explorer{
				ProjectID: "test-project",
				SummaryFields: &gcplogurl.SummaryFields{
					Fields:       []string{"trace", "protoPayload/resource"},
					Truncate:     true,
					MaxLen:       16,
					TruncateFrom: gcplogurl.TruncateFromBeginning,
				},
			},
			want: "https://console.cloud.google.com/logs/query;summaryFields=trace,protoPayload%252Fresource:true:16:beginning?project=test-project",
		},
		{
			name: "summary fields with zero value",
			obj: &gcplogurl.Explorer{
				ProjectID: "test-project",
				SummaryFields: &gcplogurl.SummaryFields{
					Fields: []string{"trace", "protoPayload/resource"},
				},
			},
			want: "https://console.cloud.google.com/logs/query;summaryFields=trace,protoPayload%252Fresource:false:32:end?project=test-project",
		},
		{
			name: "summary fields with strange value",
			obj: &gcplogurl.Explorer{
				ProjectID: "test-project",
				SummaryFields: &gcplogurl.SummaryFields{
					Fields: []string{"foo;bar", "foo bar"},
				},
			},
			want: "https://console.cloud.google.com/logs/query;summaryFields=foo%253Bbar,foo%2520bar:false:32:end?project=test-project",
			// TODO: `foo bar` should be encode to escape(`"foo bar"`), but I can't guess the rule.
			// want: "https://console.cloud.google.com/logs/query;summaryFields=foo%253Bbar,%2522foo%2520bar%2522:false:32:end?project=test-project",
		},
		{
			name: "custom fields",
			obj: &gcplogurl.Explorer{
				ProjectID:    "test-project",
				CustomFields: []string{"jsonPayload/statusDetails"},
			},
			want: "https://console.cloud.google.com/logs/query;lfeCustomFields=jsonPayload%252FstatusDetails?project=test-project",
		},
		{
			name: "custom fields - with two fields",
			obj: &gcplogurl.Explorer{
				ProjectID:    "test-project",
				CustomFields: []string{"jsonPayload/statusDetails", "jsonPayload/enforcedSecurityPolicy/name"},
			},
			want: "https://console.cloud.google.com/logs/query;lfeCustomFields=jsonPayload%252FstatusDetails,jsonPayload%252FenforcedSecurityPolicy%252Fname?project=test-project",
		},
		// all in one!
		{
			name: "all in one",
			obj: &gcplogurl.Explorer{
				ProjectID: "test-project",
				StorageScope: &gcplogurl.StorageScopeStorage{
					Src: []string{
						"projects/test-project/locations/global/buckets/_Default/views/_AllLogs",
						"projects/test-project/locations/global/buckets/_Default/views/_Default",
					},
				},
				TimeRange: &gcplogurl.SpecificTimeWithRange{
					At:    time.Date(2020, 11, 14, 4, 20, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
					Range: 2 * time.Hour,
				},
				SummaryFields: &gcplogurl.SummaryFields{
					Fields:       []string{"trace", "protoPayload/resource"},
					Truncate:     true,
					MaxLen:       16,
					TruncateFrom: gcplogurl.TruncateFromBeginning,
				},
				CustomFields: []string{"jsonPayload/statusDetails"},
			},
			want: "https://console.cloud.google.com/logs/query;lfeCustomFields=jsonPayload%252FstatusDetails;storageScope=storage,projects%2Ftest-project%2Flocations%2Fglobal%2Fbuckets%2F_Default%2Fviews%2F_AllLogs,projects%2Ftest-project%2Flocations%2Fglobal%2Fbuckets%2F_Default%2Fviews%2F_Default;summaryFields=trace,protoPayload%252Fresource:true:16:beginning;timeRange=2020-11-13T19:20:00Z%2F2020-11-13T19:20:00Z--PT2H?project=test-project",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			lp := tt.obj
			got := lp.String()
			if got != tt.want {
				t.Errorf("Build() got = %v, want %v", got, tt.want)
				u1, err := url.Parse(got)
				if err != nil {
					t.Fatal(err)
				}
				u2, err := url.Parse(tt.want)
				if err != nil {
					t.Fatal(err)
				}
				t.Log("got ", u1.Path, u1.RawPath)
				t.Log("want", u2.Path, u2.RawPath)
			}
		})
	}
}
