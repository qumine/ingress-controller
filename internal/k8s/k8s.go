package k8s

import (
	"context"
	"net"
	"strconv"

	"github.com/quhive/qumine-ingress/internal/server"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// AnnotationHostname is the kubernetes annotation for the hostname to use for the ingress
	AnnotationHostname = "qumine-ingress.quhive.io/hostname"
	// AnnotationPortname is the kubernetes annotation for the name of the port to use
	AnnotationPortname = "qumine-ingress.quhive.io/portname"
)

// K8S is a watcher for kubernetes
type K8S struct {
	// Status is the current status of the K8S watcher.
	Status string

	config *rest.Config
	stop   chan struct{}
}

// NewK8S creates a new k8s instance
func NewK8S(kubeConfig string) *K8S {
	k8s := &K8S{}
	k8s.stop = make(chan struct{}, 1)

	if kubeConfig != "" {
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			logrus.WithError(err).Fatal("unable to load kube-config")
		}

		k8s.config = config
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			logrus.WithError(err).Fatal("unable to load in-cluster config")
		}

		k8s.config = config
	}

	return k8s
}

// Start the K8S
func (k8s *K8S) Start(context context.Context) {
	defer k8s.close()
	logrus.Info("starting k8s...")

	clientset, err := kubernetes.NewForConfig(k8s.config)
	if err != nil {
		logrus.WithError(err).Fatal("unable to create kubernetes clientset")
	}

	watchlist := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		string(v1.ResourceServices),
		v1.NamespaceAll,
		fields.Everything(),
	)

	_, controller := cache.NewInformer(
		watchlist,
		&v1.Service{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				service := extractRoutableService(obj)
				if service != nil {
					if service.hostname != "" {
						metricsRoutes.Inc()
						server.CreateRoute(service.hostname, service.backend)
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				service := extractRoutableService(obj)
				if service != nil {
					if service.hostname != "" {
						metricsRoutes.Dec()
						server.DeleteRoute(service.hostname)
					}
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldService := extractRoutableService(oldObj)
				newService := extractRoutableService(newObj)
				if oldService != nil && newService != nil {
					if oldService.hostname != "" && newService.hostname != "" {
						server.DeleteRoute(oldService.hostname)
						server.CreateRoute(newService.hostname, newService.backend)
					}
				}
			},
		},
	)

	go controller.Run(k8s.stop)
	k8s.Status = "up"
	for {
		select {
		case <-context.Done():
			k8s.Status = "down"
			return
		}
	}
}

func (k8s *K8S) close() {
	k8s.stop <- struct{}{}
}

// Service represents the service definition from kubernetes.
type Service struct {
	hostname string
	backend  string
}

func extractRoutableService(obj interface{}) *Service {
	service, ok := obj.(*v1.Service)
	if !ok {
		return nil
	}

	portname := "minecraft"
	if p, exists := service.Annotations[AnnotationHostname]; exists {
		portname = p
	}

	if hostname, exists := service.Annotations[AnnotationHostname]; exists {
		return buildDetails(service, hostname, portname)
	}

	return nil
}

func buildDetails(service *v1.Service, hostname string, portname string) *Service {
	clusterIP := service.Spec.ClusterIP
	port := "25565"
	for _, p := range service.Spec.Ports {
		if p.Name == portname {
			port = strconv.Itoa(int(p.Port))
		}
	}
	rs := &Service{
		hostname: hostname,
		backend:  net.JoinHostPort(clusterIP, port),
	}
	return rs
}
