package k8s

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricsRoutes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "qumine_ingress_routes",
			Help: "The amount of registered routes",
		},
	)
)

func init() {
	prometheus.MustRegister(metricsRoutes)
}
