package gcplogurl

import (
	"io"
	"net/url"
	"time"

	"github.com/rickb777/date/period"
)

// TimeRange means when to display the logs.
type TimeRange interface {
	isTimeRange()
	marshalURL(w io.Writer)
}

var _ TimeRange = (*RecentRange)(nil)
var _ TimeRange = (*SpecificTimeBetween)(nil)
var _ TimeRange = (*SpecificTimeWithRange)(nil)

// RecentRange provides "** seconds/minutes/hours/days ago".
type RecentRange struct {
	Last time.Duration
}

func (t *RecentRange) isTimeRange() {}

func (t *RecentRange) marshalURL(w io.Writer) {
	_, _ = w.Write([]byte(";timeRange="))
	p, _ := period.NewOf(t.Last)
	_, _ = w.Write([]byte(p.String()))
}

// SpecificTimeBetween pvovides custom range.
type SpecificTimeBetween struct {
	From time.Time
	To   time.Time
}

func (t *SpecificTimeBetween) isTimeRange() {}

func (t *SpecificTimeBetween) marshalURL(w io.Writer) {
	_, _ = w.Write([]byte(";timeRange="))
	if v := t.From; !v.IsZero() {
		_, _ = w.Write([]byte(v.In(time.UTC).Format(time.RFC3339Nano)))
	}
	_, _ = w.Write([]byte(url.QueryEscape("/")))
	if v := t.To; !v.IsZero() {
		_, _ = w.Write([]byte(v.In(time.UTC).Format(time.RFC3339Nano)))
	}
}

// SpecificTimeWithRange provides jump tp time.
type SpecificTimeWithRange struct {
	At    time.Time
	Range time.Duration
}

func (t *SpecificTimeWithRange) isTimeRange() {}

func (t *SpecificTimeWithRange) marshalURL(w io.Writer) {
	_, _ = w.Write([]byte(";timeRange="))
	_, _ = w.Write([]byte(t.At.In(time.UTC).Format(time.RFC3339Nano)))
	_, _ = w.Write([]byte(url.QueryEscape("/")))
	_, _ = w.Write([]byte(t.At.In(time.UTC).Format(time.RFC3339Nano)))
	_, _ = w.Write([]byte(url.QueryEscape("--")))
	p, _ := period.NewOf(t.Range)
	_, _ = w.Write([]byte(p.String()))
}
