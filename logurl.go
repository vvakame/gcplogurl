package gcplogurl

import (
	"bytes"
	"net/url"
)

const logBaseURL = "https://console.cloud.google.com/logs/query"

// Explorer is a map of GCP Cloud Logging Log Explorer.
type Explorer struct {
	BaseURL   string
	ProjectID string
	// Query expression for logs.
	// https://cloud.google.com/logging/docs/view/logging-query-language
	Query string
	// StorageScope for refine scope.
	StorageScope StorageScope
	// TimeRange for filter logs.
	TimeRange TimeRange
	// SummaryFields for manage summary fields.
	SummaryFields *SummaryFields
}

// String returns represent of Explorer URL.
func (ex *Explorer) String() string {
	baseURL := ex.BaseURL
	if baseURL == "" {
		baseURL = logBaseURL
	}
	w := bytes.NewBufferString(baseURL)
	if v := ex.Query; v != "" {
		w.Write([]byte(";query="))
		w.Write([]byte(url.QueryEscape(v)))
	}
	if v := ex.StorageScope; v != nil {
		v.marshalURL(w)
	}
	if v := ex.TimeRange; v != nil {
		v.marshalURL(w)
	}
	if v := ex.SummaryFields; v != nil {
		v.marshalURL(w)
	}
	if v := ex.ProjectID; v != "" {
		w.WriteString("?project=")
		w.WriteString(url.QueryEscape(v))
	}

	return w.String()
}
