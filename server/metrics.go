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
		[]string{"hostname"},
	)
	metricsConnectionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qumine_ingress_connections_total",
			Help: "The total amount of connections",
		},
		[]string{"hostname"},
	)
	metricsBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qumine_ingress_bytes_total",
			Help: "The total bytes transmitted",
		},
		[]string{"direction"},
	)
)

func init() {
	prometheus.MustRegister(metricsConnections)
	prometheus.MustRegister(metricsConnectionsTotal)
	prometheus.MustRegister(metricsBytes)
}
