package config

import (
	"net"
	"strconv"

	"github.com/spf13/pflag"
)

var ingressOptions IngressOptions

type IngressOptions struct {
	Host string
	Port int
}

func GetIngressFlagSet() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&apiOptions.Host, "host", "", "Host for the API server to listen on (default: 0.0.0.0)")
	flagSet.IntVar(&apiOptions.Port, "port", 8080, "Port for the API server to listen on (default: 8080)")
	return flagSet
}

func GetIngressOptions() IngressOptions {
	return ingressOptions
}

func (ingressOptions *IngressOptions) GetAddress() string {
	return net.JoinHostPort(ingressOptions.Host, strconv.Itoa(ingressOptions.Port))
}
