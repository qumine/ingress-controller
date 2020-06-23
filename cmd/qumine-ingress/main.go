package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/quhive/qumine-ingress/internal/api"
	"github.com/quhive/qumine-ingress/internal/k8s"
	"github.com/quhive/qumine-ingress/internal/server"

	"github.com/sirupsen/logrus"
)

var (
	version = "dev"
	commit  = "none"
	date    = "uknown"
)

var (
	helpFlag    bool
	versionFlag bool
	debugFlag   bool
	traceFlag   bool

	kubeConfig string

	host string
	port int
)

func init() {
	flag.BoolVar(&helpFlag, "help", false, "Show this page")
	flag.BoolVar(&versionFlag, "version", false, "Show the current version")
	flag.BoolVar(&debugFlag, "debug", false, "Enable debugging log level")
	flag.BoolVar(&traceFlag, "trace", false, "Enable trace log level")

	flag.StringVar(&kubeConfig, "kube-config", "", "Path of the kube config file to use")

	flag.StringVar(&host, "host", "0.0.0.0", "Address the server will listen on")
	flag.IntVar(&port, "port", 25565, "Port the server will listen on")
	flag.Parse()
}

func main() {
	if helpFlag {
		showHelp()
	}

	if versionFlag {
		showVersion()
	}

	if debugFlag {
		enableDebug()
	}
	if debugFlag {
		enableTrace()
	}

	api := api.NewAPI()
	k8s := k8s.NewK8S(kubeConfig)
	server := server.NewServer(host, port)

	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go k8s.Start(context)
	go api.Start(context, k8s, server)
	go server.Start(context)

	<-c
}

func showHelp() {
	flag.Usage()
	os.Exit(0)
}

func showVersion() {
	fmt.Printf("%v, commit %v, built at %v", version, commit, date)
	os.Exit(0)
}

func enableDebug() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debug("debugging enabled")
}

func enableTrace() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.Debug("tracing enabled")
}
