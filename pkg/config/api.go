package config

import (
	"net"
	"strconv"

	"github.com/spf13/pflag"
)

var apiOptions APIOptions

type APIOptions struct {
	Host string
	Port int
}

func GetAPIFlagSet() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&apiOptions.Host, "api-host", "0.0.0.0", "Host for the API server to listen on")
	flagSet.IntVar(&apiOptions.Port, "api-port", 8080, "Port for the API server to listen on")
	return flagSet
}

func GetAPIOptions() APIOptions {
	return apiOptions
}

func (aiOptions *APIOptions) GetAddress() string {
	return net.JoinHostPort(aiOptions.Host, strconv.Itoa(apiOptions.Port))
}
