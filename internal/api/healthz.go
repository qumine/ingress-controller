package api

import "net/http"

func (api *API) getHealthz(writer http.ResponseWriter, request *http.Request) {
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
