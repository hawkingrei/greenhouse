/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// DefaultFieldsFormatter wraps another logrus.Formatter, injecting
// DefaultFields into each Format() call, existing fields are preserved
// if they have the same key
type DefaultFieldsFormatter struct {
	WrappedFormatter logrus.Formatter
	DefaultFields    logrus.Fields
}

// NewDefaultFieldsFormatter returns a DefaultFieldsFormatter,
// if wrappedFormatter is nil &logrus.JSONFormatter{} will be used instead
func NewDefaultFieldsFormatter(
	wrappedFormatter logrus.Formatter, defaultFields logrus.Fields,
) *DefaultFieldsFormatter {
	res := &DefaultFieldsFormatter{
		WrappedFormatter: wrappedFormatter,
		DefaultFields:    defaultFields,
	}
	if res.WrappedFormatter == nil {
		res.WrappedFormatter = &logrus.JSONFormatter{}
	}
	return res
}

// Format implements logrus.Formatter's Format. We allocate a a new Fields
// map in order to not modify the caller's Entry, as that is not a thread
// safe operation.
func (d *DefaultFieldsFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+len(d.DefaultFields))
	for k, v := range d.DefaultFields {
		data[k] = v
	}
	for k, v := range entry.Data {
		data[k] = v
	}
	return d.WrappedFormatter.Format(&logrus.Entry{
		Logger:  entry.Logger,
		Data:    data,
		Time:    entry.Time,
		Level:   entry.Level,
		Message: entry.Message,
	})
}

// prometheusMetrics are served by /prometheus on the metrics port
type prometheusMetrics struct {
	DiskFree             prometheus.Gauge
	DiskUsed             prometheus.Gauge
	DiskTotal            prometheus.Gauge
	FilesEvicted         prometheus.Counter
	ActionCacheHits      prometheus.Counter
	CASHits              prometheus.Counter
	ActionCacheMisses    prometheus.Counter
	CASMisses            prometheus.Counter
	LastEvictedAccessAge prometheus.Gauge
}

func initMetrics() *prometheusMetrics {
	metrics := &prometheusMetrics{
		DiskFree: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "bazel_cache_disk_free",
			Help: "Free gb on bazel cache disk",
		}),
		DiskUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "bazel_cache_disk_used",
			Help: "Used gb on bazel cache disk",
		}),
		DiskTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "bazel_cache_disk_total",
			Help: "Total gb on bazel cache disk",
		}),
		FilesEvicted: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "bazel_cache_evicted_files",
			Help: "number of files evicted since last server start",
		}),
		ActionCacheHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "bazel_cache_cas_hits",
			Help: "Approximate number of Action Cache hits since last server start",
		}),
		CASHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "bazel_cache_action_hits",
			Help: "Approximate number of Content Addressed Storage cache hits since last server start",
		}),
		ActionCacheMisses: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "bazel_cache_action_misses",
			Help: "Approximate number of Content Addressed Storage cache misses since last server start",
		}),
		CASMisses: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "bazel_cache_cas_misses",
			Help: "Approximate number of Content Addressed Storage cache misses since last server start",
		}),
		LastEvictedAccessAge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "bazel_cache_last_evicted_access_age",
			Help: "Hours since last access of most recently evicted file (at eviction time)",
		}),
	}
	prometheus.MustRegister(metrics.DiskFree)
	prometheus.MustRegister(metrics.DiskUsed)
	prometheus.MustRegister(metrics.DiskTotal)
	prometheus.MustRegister(metrics.FilesEvicted)
	prometheus.MustRegister(metrics.ActionCacheHits)
	prometheus.MustRegister(metrics.CASHits)
	prometheus.MustRegister(metrics.ActionCacheMisses)
	prometheus.MustRegister(metrics.CASMisses)
	prometheus.MustRegister(metrics.LastEvictedAccessAge)
	return metrics
}
