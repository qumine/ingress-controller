package k8s

import (
	"net"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/quhive/qumine-ingress/internal/metrics"
	"github.com/quhive/qumine-ingress/internal/routing"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

func onAdd(obj interface{}) {
	service, ok := obj.(*v1.Service)
	if !ok {
		metrics.ErrorsTotal.With(prometheus.Labels{"error": "InternalError"}).Inc()
		return
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
			return
		}
	}
	metrics.ErrorsTotal.With(prometheus.Labels{"error": "NoMatchingPort"}).Inc()
}

func onDelete(obj interface{}) {
	service, ok := obj.(*v1.Service)
	if !ok {
		metrics.ErrorsTotal.With(prometheus.Labels{"error": "InternalError"}).Inc()
		return
	}

	routing.Remove(string(service.UID))
}

func onUpdate(oldObj, newObj interface{}) {
	oldService, ok := newObj.(*v1.Service)
	if !ok {
		metrics.ErrorsTotal.With(prometheus.Labels{"error": "InternalError"}).Inc()
		return
	}
	newService, ok := newObj.(*v1.Service)
	if !ok {
		metrics.ErrorsTotal.With(prometheus.Labels{"error": "InternalError"}).Inc()
		return
	}

	portname := "minecraft"
	if p, exists := newService.Annotations[AnnotationPortname]; exists {
		portname = p
	}
	hostname := "localhost"
	if h, exists := newService.Annotations[AnnotationHostname]; exists {
		hostname = h
	}
	logrus.WithField("portname", portname).WithField("hostname", hostname).Debug("add route")

	for _, p := range newService.Spec.Ports {
		if p.Name == portname {
			routing.Update(string(oldService.UID), routing.NewRoute(hostname, net.JoinHostPort(newService.Spec.ClusterIP, strconv.Itoa(int(p.Port)))))
			return
		}
	}
	metrics.ErrorsTotal.With(prometheus.Labels{"error": "NoMatchingPort"}).Inc()
}
