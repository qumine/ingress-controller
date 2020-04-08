package api

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/quhive/qumine-ingress/server"
	"github.com/sirupsen/logrus"
)

// API represents the api server
type API struct {
	server *http.Server
	router *mux.Router
}

// NewAPI creates a new api instance with the given host and port
func NewAPI(host string, port int) *API {
	router := mux.NewRouter()
	router.Use(metricsMiddleware)
	router.Use(loggingMiddleware)
	router.Path("/routes").Methods("GET").HandlerFunc(getRoutes)
	router.Path("/metrics").Methods("GET").Handler(promhttp.Handler())

	return &API{
		server: &http.Server{
			Addr:    net.JoinHostPort(host, strconv.Itoa(port)),
			Handler: router,
		},
		router: router,
	}
}

// Start the Api
func (api *API) Start(context context.Context) {
	defer api.server.Close()
	logrus.WithFields(logrus.Fields{
		"addr": api.server.Addr,
	}).Info("starting api...")

	go logrus.WithError(api.server.ListenAndServe()).WithFields(logrus.Fields{
		"addr": api.server.Addr,
	}).Fatal("api failed to start")

	for {
		select {
		case <-context.Done():
			return
		}
	}
}

func getRoutes(writer http.ResponseWriter, request *http.Request) {
	mappings := server.GetMappings()
	bytes, err := json.Marshal(mappings)
	if err != nil {
		logrus.WithError(err).Error("marchaling mappings failed")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(bytes)
}
