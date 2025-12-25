package lib

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMetricsHandlerWithRegistry(
	mux *http.ServeMux,
	registry *prometheus.Registry,
) {
	handler := promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{},
	)
	mux.Handle("/metrics", handler)
}
