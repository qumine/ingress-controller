package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/quhive/qumine-ingress/internal/k8s"
	"github.com/quhive/qumine-ingress/internal/server"
	"github.com/sirupsen/logrus"
)

var (
	k *k8s.K8S
	s *server.Server
)

// API represents the api server
type API struct {
	httpServer *http.Server
	router     *mux.Router
}

// NewAPI creates a new api instance with the given host and port
func NewAPI() *API {
	router := mux.NewRouter()
	router.Path("/healthz").Methods("GET").HandlerFunc(getHealthz)
	router.Path("/metrics").Methods("GET").Handler(promhttp.Handler())
	return &API{
		httpServer: &http.Server{
			Addr:    "0.0.0.0:8080",
			Handler: router,
		},
		router: router,
	}
}

// Start the Api
func (api *API) Start(context context.Context, k8s *k8s.K8S, server *server.Server) {
	defer api.httpServer.Close()
	logrus.WithFields(logrus.Fields{
		"addr": api.httpServer.Addr,
	}).Info("starting api...")

	k = k8s
	s = server

	go logrus.WithError(api.httpServer.ListenAndServe()).WithFields(logrus.Fields{
		"addr": api.httpServer.Addr,
	}).Fatal("api failed to start")

	for {
		select {
		case <-context.Done():
			return
		}
	}
}

func getHealthz(writer http.ResponseWriter, request *http.Request) {
	details := make(map[string]string)
	details["k8s"] = k.Status
	details["server"] = s.Status

	status := "up"
	for key := range details {
		switch details[key] {
		case "up":
			continue
		default:
			status = "down"
			break
		}
	}

	bytes, err := json.Marshal(&healthz{
		Status:  status,
		Details: details,
	})
	if err != nil {
		logrus.WithError(err).Error("marchaling healthz failed")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if status == "down" {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}

	writer.Write(bytes)
}

type healthz struct {
	Status  string            `json:"status"`
	Details map[string]string `json:"details"`
}
