package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Routes represents the metrics for the amount of registered routes
	Routes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "qumine_ingress_routes",
			Help: "The amount of registered routes",
		},
	)
	// Connections represents the metrics for the amount of active connections
	Connections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "qumine_ingress_connections",
			Help: "The amount of active connections",
		},
		[]string{"route"},
	)
	// ErrorsTotal represents the metrics for the amount of total errors
	ErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qumine_ingress_errors_total",
			Help: "The total error count",
		},
		[]string{"error"},
	)
	// BytesTotal represents the metrics for the amount of total bytes transmitted
	BytesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qumine_ingress_bytes_total",
			Help: "The total bytes transmitted",
		},
		[]string{"direction", "route"},
	)
)

func init() {
	prometheus.MustRegister(Routes)
	prometheus.MustRegister(Connections)
	prometheus.MustRegister(ErrorsTotal)
	prometheus.MustRegister(BytesTotal)
}
