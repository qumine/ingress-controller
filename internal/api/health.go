package api

import "net/http"

func (api *API) healthLive(writer http.ResponseWriter, request *http.Request) {
	details := make(map[string]string)
	details["k8s"] = api.k8s.Status
	details["server"] = api.ing.Status

	if api.k8s.Status == "up" && api.ing.Status == "up" {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte{})
	} else {
		writer.WriteHeader(http.StatusServiceUnavailable)
		writer.Write([]byte{})
	}
}

func (api *API) healthReady(writer http.ResponseWriter, request *http.Request) {
	details := make(map[string]string)
	details["k8s"] = api.k8s.Status
	details["server"] = api.ing.Status

	if api.k8s.Status == "up" && api.ing.Status == "up" {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte{})
	} else {
		writer.WriteHeader(http.StatusServiceUnavailable)
		writer.Write([]byte{})
	}
}
