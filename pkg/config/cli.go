package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	debug bool
	trace bool
)

type CliOptions struct {
	LogLevel logrus.Level
}

func GetCliFlagSet() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.BoolVarP(&debug, "debug", "d", false, "Debug logging (default: false)")
	flagSet.BoolVar(&trace, "trace", false, "Trace logging (default: false)")
	return flagSet
}

func GetCliOptions() CliOptions {
	logLevel := logrus.InfoLevel
	if debug {
		logLevel = logrus.DebugLevel
	}
	if trace {
		logLevel = logrus.TraceLevel
	}

	return CliOptions{
		LogLevel: logLevel,
	}
}
