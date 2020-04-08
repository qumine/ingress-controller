package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		logrus.WithFields(logrus.Fields{
			"client": request.RemoteAddr,
			"method": request.Method,
			"url":    request.URL,
		}).Info("inbound api request")

		next.ServeHTTP(writer, request)
	})
}
