package api

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qumine/ingress-controller/internal/ingress"
	"github.com/qumine/ingress-controller/internal/k8s"
	"github.com/qumine/ingress-controller/pkg/config"
	"github.com/sirupsen/logrus"
)

// API represents the api server
type API struct {
	addr string
	k8s  *k8s.K8S
	ing  *ingress.Ingress
}

// NewAPI creates a new api instance with the given host and port
func NewAPI(apiOptions config.APIOptions, k8s *k8s.K8S, ing *ingress.Ingress) *API {
	return &API{
		addr: apiOptions.GetAddress(),
		k8s:  k8s,
		ing:  ing,
	}
}

// Start the Api
func (api *API) Start(context context.Context) {
	logrus.WithFields(logrus.Fields{
		"addr": api.addr,
	}).Debug("Starting API")

	go api.startHttpServer()

	logrus.WithFields(logrus.Fields{
		"addr": api.addr,
	}).Info("Started API")
	for {
		select {
		case <-context.Done():
			return
		}
	}
}

func (api *API) startHttpServer() {
	router := http.NewServeMux()
	router.HandleFunc("/healthz", api.getHealthz)
	router.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:    api.addr,
		Handler: router,
	}
	defer httpServer.Close()
	if err := httpServer.ListenAndServe(); err != nil {
		logrus.WithFields(logrus.Fields{
			"addr": api.addr,
		}).Fatal("Failed to start API")
	}
}
