package api

import (
	"context"
	"net/http"
	"sync"

	"github.com/qumine/ingress-controller/internal/ingress"
	"github.com/qumine/ingress-controller/internal/k8s"
	"github.com/qumine/ingress-controller/pkg/config"
	"github.com/sirupsen/logrus"
)

// API represents the api server
type API struct {
	k8s *k8s.K8S
	ing *ingress.Ingress

	httpServer *http.Server
}

// NewAPI creates a new api instance with the given host and port
func NewAPI(apiOptions config.APIOptions, k8s *k8s.K8S, ing *ingress.Ingress) *API {
	r := http.NewServeMux()
	api := &API{
		k8s: k8s,
		ing: ing,

		httpServer: &http.Server{
			Addr:    apiOptions.GetAddress(),
			Handler: r,
		},
	}
	r.HandleFunc("/health/live", api.healthLive)
	r.HandleFunc("/health/ready", api.healthReady)

	return api
}

// Start the Api
func (api *API) Start(context context.Context, wg *sync.WaitGroup) {
	defer api.Stop(wg)
	logrus.WithFields(logrus.Fields{
		"addr": api.httpServer.Addr,
	}).Debug("Starting API")

	go func() {
		wg.Add(1)
		if err := api.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithFields(logrus.Fields{
				"addr": api.httpServer.Addr,
			}).Fatal("Failed to start API")
		}
	}()

	logrus.WithFields(logrus.Fields{
		"addr": api.httpServer.Addr,
	}).Info("Started API")
	for {
		select {
		case <-context.Done():
			return
		}
	}
}

// Stop the api
func (a *API) Stop(wg *sync.WaitGroup) {
	logrus.WithFields(logrus.Fields{
		"addr": a.httpServer.Addr,
	}).Debug("Stopping API")

	if err := a.httpServer.Close(); err != nil {
		logrus.WithFields(logrus.Fields{
			"addr": a.httpServer.Addr,
		}).Error("Failed to stop API")
	}

	wg.Done()
	logrus.WithFields(logrus.Fields{
		"addr": a.httpServer.Addr,
	}).Info("Stopped API")
}
