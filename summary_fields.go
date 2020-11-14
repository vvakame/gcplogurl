package gcplogurl

import (
	"io"
	"net/url"
	"strconv"
)

// TruncateFrom provides option for truncate field name.
type TruncateFrom string

const (
	// TruncateFromBeginning about field name.
	TruncateFromBeginning TruncateFrom = "beginning"
	// TruncateFromEnd about field name.
	TruncateFromEnd TruncateFrom = "end"
)

// SummaryFields provides configurations for log's summary fields.
type SummaryFields struct {
	Fields       []string
	Truncate     bool
	MaxLen       int
	TruncateFrom TruncateFrom
}

func (sf *SummaryFields) marshalURL(w io.Writer) {
	_, _ = w.Write([]byte(";summaryFields="))
	for idx, f := range sf.Fields {
		if idx != 0 {
			_, _ = w.Write([]byte(","))
		}
		_, _ = w.Write([]byte(url.QueryEscape(f)))
	}
	_, _ = w.Write([]byte(":"))
	_, _ = w.Write([]byte(strconv.FormatBool(sf.Truncate)))
	_, _ = w.Write([]byte(":"))
	ml := sf.MaxLen
	if ml == 0 {
		ml = 32
	}
	_, _ = w.Write([]byte(strconv.Itoa(ml)))
	_, _ = w.Write([]byte(":"))
	tf := sf.TruncateFrom
	if tf == "" {
		tf = TruncateFromEnd
	}
	_, _ = w.Write([]byte(tf))
}
