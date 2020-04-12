package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricsConnectionsTotal = prometheus.NewCounterVec(
		prometheus.GaugeOpts{
			Name: "qumine_ingress_connections_total",
			Help: "The total connections count",
		},
		[]string{"route"},
	)
	metricsErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qumine_ingress_errors_total",
			Help: "The total error count",
		},
		[]string{"error"},
	)
	metricsBytesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qumine_ingress_bytes_total",
			Help: "The total bytes transmitted",
		},
		[]string{"direction", "route"},
	)
)

func init() {
	prometheus.MustRegister(metricsConnectionsTotal)
	prometheus.MustRegister(metricsErrorsTotal)
	prometheus.MustRegister(metricsBytesTotal)
}
