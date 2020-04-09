package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricsConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "qumine_ingress_connections",
			Help: "The amount of active connections",
		},
		[]string{"route"},
	)
	metricsErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qumine_ingress_errors_total",
			Help: "The total error count",
		},
		[]string{"error"},
	)
	metricsBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qumine_ingress_bytes_total",
			Help: "The total bytes transmitted",
		},
		[]string{"direction", "route"},
	)
)

func init() {
	prometheus.MustRegister(metricsConnections)
	prometheus.MustRegister(metricsBytes)
}
