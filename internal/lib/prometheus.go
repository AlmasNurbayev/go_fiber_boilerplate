package lib

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type PrometheusType struct {
	Registry     *prometheus.Registry
	CounterVec   *prometheus.CounterVec
	HistogramVec *prometheus.HistogramVec
}

func NewPromRegistry(log *slog.Logger) PrometheusType {
	registry := prometheus.NewRegistry()
	httpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_milliseconds",
			Help:    "Duration of HTTP requests in milliseconds",
			Buckets: []float64{1, 10, 50, 100, 200, 500, 1000}, // Бакеты аналогичны JS
		},
		[]string{"method", "route", "statusCode", "originalUrl"},
	)
	httpRequestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "statusCode", "originalUrl"},
	)

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		httpRequestDuration,
		httpRequestCounter,
	)
	log.Info("init prometheus registry")

	return PrometheusType{
		Registry:     registry,
		CounterVec:   httpRequestCounter,
		HistogramVec: httpRequestDuration,
	}
}
