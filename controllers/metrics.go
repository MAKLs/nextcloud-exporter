package controllers

import (
	"net/http"

	"github.com/MAKLs/nextcloud-exporter/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsController serves the Prometheus metrics collected by the exporter.
func MetricsController() http.Handler {
	return promhttp.HandlerFor(metrics.ExporterRegistry, promhttp.HandlerOpts{})
}
