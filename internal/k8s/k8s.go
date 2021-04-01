package k8s

import (
	"context"

	"github.com/qumine/ingress-controller/pkg/config"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// AnnotationHostname is the kubernetes annotation for the hostname to use for the ingress
	AnnotationHostname = "ingress.qumine.io/hostname"
	// AnnotationPortname is the kubernetes annotation for the name of the port to use
	AnnotationPortname = "ingress.qumine.io/portname"
)

// K8S is a watcher for kubernetes
type K8S struct {
	// Status is the current status of the K8S watcher.
	Status string

	kubeconfig string
	stop       chan struct{}
}

// NewK8S creates a new k8s instance
func NewK8S(k8sOptions config.K8SOptions) *K8S {
	return &K8S{
		kubeconfig: k8sOptions.KubeConfig,
		stop:       make(chan struct{}, 1),
	}
}

// Start the K8S
func (k8s *K8S) Start(context context.Context) {
	defer k8s.close()
	logrus.WithFields(logrus.Fields{
		"kubeconfig": k8s.kubeconfig,
	}).Debug("Starting K8S")

	config, err := clientcmd.BuildConfigFromFlags("", k8s.kubeconfig)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"kubeconfig": k8s.kubeconfig,
		}).Fatal("Failed to start K8S")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"kubeconfig": k8s.kubeconfig,
		}).Fatal("Failed to start K8S")
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
			AddFunc:    onAdd,
			DeleteFunc: onDelete,
			UpdateFunc: onUpdate,
		},
	)

	go controller.Run(k8s.stop)
	k8s.Status = "up"
	logrus.WithFields(logrus.Fields{
		"kubeconfig": k8s.kubeconfig,
	}).Info("Started K8S")
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
