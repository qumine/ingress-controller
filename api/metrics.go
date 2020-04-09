package api

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricsAPIRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qumine_ingress_api_requests_total",
			Help: "The total number of api requests",
		},
		[]string{"path"},
	)
	metricsAPITime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "qumine_ingress_api_time",
			Help:    "The response time of the api",
			Buckets: prometheus.LinearBuckets(0.01, 0.01, 10),
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(metricsAPIRequestsTotal)
	prometheus.MustRegister(metricsAPITime)
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		metricsAPIRequestsTotal.With(prometheus.Labels{"path": request.URL.Path}).Inc()

		start := time.Now()
		next.ServeHTTP(writer, request)
		metricsAPITime.With(prometheus.Labels{"path": request.URL.Path}).Observe(float64(time.Since(start).Milliseconds()))
	})
}
