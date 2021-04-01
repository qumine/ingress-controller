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
	flagSet.StringVar(&ingressOptions.Host, "host", "", "Host for the API server to listen on (default: 0.0.0.0)")
	flagSet.IntVar(&ingressOptions.Port, "port", 25565, "Port for the API server to listen on (default: 25565)")
	return flagSet
}

func GetIngressOptions() IngressOptions {
	return ingressOptions
}

func (ingressOptions *IngressOptions) GetAddress() string {
	return net.JoinHostPort(ingressOptions.Host, strconv.Itoa(ingressOptions.Port))
}
