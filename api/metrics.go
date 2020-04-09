package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricsAPIRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qumine_ingress_api_requests_total",
			Help: "The total number of api requests",
		},
		[]string{"path"},
	)
	metricsAPITime = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "qumine_ingress_api_time",
		Help:    "The response time of the api",
		Buckets: prometheus.LinearBuckets(0.01, 0.01, 10),
	})
)

func init() {
	prometheus.MustRegister(metricsAPIRequestsTotal)
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		timer := prometheus.NewTimer(metricsAPITime)
		defer timer.ObserveDuration()

		metricsAPIRequestsTotal.With(prometheus.Labels{"path": request.URL.Path}).Inc()
		next.ServeHTTP(writer, request)
	})
}
