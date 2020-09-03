package k8s

import (
	"errors"
	"net"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/quhive/qumine-ingress/internal/metrics"
	"github.com/quhive/qumine-ingress/internal/routing"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

func onAdd(obj interface{}) {
	if err := add(obj); err != nil {
		return
	}
	metrics.Routes.Inc()
	return
}

func onDelete(obj interface{}) {
	if err := remove(obj); err != nil {
		return
	}
	metrics.Routes.Dec()
}

func onUpdate(oldObj, newObj interface{}) {
	if err := remove(oldObj); err != nil {
		return
	}
	if err := add(newObj); err != nil {
		return
	}
}

func add(obj interface{}) error {
	service, ok := obj.(*v1.Service)
	if !ok {
		metrics.ErrorsTotal.With(prometheus.Labels{"error": "InternalError"}).Inc()
		return errors.New("unable to convert")
	}

	portname := "minecraft"
	if p, exists := service.Annotations[AnnotationPortname]; exists {
		portname = p
	}
	hostname := "localhost"
	if h, exists := service.Annotations[AnnotationHostname]; exists {
		hostname = h
	}
	logrus.WithField("portname", portname).WithField("hostname", hostname).Debug("add route")

	for _, p := range service.Spec.Ports {
		if p.Name == portname {
			routing.Add(string(service.UID), routing.NewRoute(hostname, net.JoinHostPort(service.Spec.ClusterIP, strconv.Itoa(int(p.Port)))))
			return nil
		}
	}
	metrics.ErrorsTotal.With(prometheus.Labels{"error": "NoMatchingPort"}).Inc()
	return errors.New("No matching port found")
}

func remove(obj interface{}) error {
	service, ok := obj.(*v1.Service)
	if !ok {
		metrics.ErrorsTotal.With(prometheus.Labels{"error": "InternalError"}).Inc()
		return errors.New("unable to convert")
	}

	routing.Remove(string(service.UID))
	return nil
}
