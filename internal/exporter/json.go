// Copyright 2015 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package exporter

import (
	"encoding/json"
	"expvar"
	"github.com/google/mtail/internal/metrics"
	"github.com/google/mtail/internal/metrics/datum"
	"net/http"

	"github.com/golang/glog"
)

type MetricResult struct {
	Name           string
	Program        string
	Kind           metrics.Kind
	Type           metrics.Type
	Hidden         bool          `json:",omitempty"`
	Keys           []string      `json:",omitempty"`
	LabelValues    []*metrics.LabelValue `json:",omitempty"`
	Source         string        `json:",omitempty"`
	Buckets        []datum.Range `json:",omitempty"`
	Limit          int64         `json:",omitempty"`
}

var exportJSONErrors = expvar.NewInt("exporter_json_errors")

// HandleJSON exports the metrics in JSON format via HTTP.
func (e *Exporter) HandleJSON(w http.ResponseWriter, r *http.Request) {

	results := make([]*MetricResult, 0)
	for _,ms := range e.store.Metrics {
		for _,m := range ms {
			rs := &MetricResult{m.Name,m.Program,m.Kind,m.Type,m.Hidden,m.Keys,m.LabelValues(),m.Source,m.Buckets,m.Limit}
			results = append(results, rs)
		}
	}

	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		exportJSONErrors.Add(1)
		glog.Info("error marshalling metrics into json:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	if _, err := w.Write(b); err != nil {
		glog.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
