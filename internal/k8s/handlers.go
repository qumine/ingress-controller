package k8s

import (
	"net"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/qumine/ingress-controller/internal/metrics"
	"github.com/qumine/ingress-controller/internal/routing"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

func onAdd(obj interface{}) {
	service, ok := obj.(*v1.Service)
	if !ok {
		metrics.ErrorsTotal.With(prometheus.Labels{"error": "InternalError"}).Inc()
		return
	}

	hostname := "localhost"
	if h, exists := service.Annotations[AnnotationHostname]; exists {
		hostname = h
	} else {
		logrus.WithFields(logrus.Fields{
			"service": service,
		}).Tracef("Adding service skipped, %s annotation not present", AnnotationHostname)
		return
	}

	portname := "minecraft"
	if p, exists := service.Annotations[AnnotationPortname]; exists {
		portname = p
	}
	logrus.WithFields(logrus.Fields{
		"hostname": hostname,
		"portname": portname,
	}).Debug("Adding route")

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

	if _, exists := service.Annotations[AnnotationHostname]; !exists {
		logrus.WithFields(logrus.Fields{
			"service": service,
		}).Tracef("Deleting service skipped, %s annotation not present", AnnotationHostname)
		return
	}

	routing.Remove(string(service.UID))
}

func onUpdate(oldObj interface{}, newObj interface{}) {
	onDelete(oldObj)
	onAdd(newObj)
}
