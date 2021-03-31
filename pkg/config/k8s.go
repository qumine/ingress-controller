package config

import "github.com/spf13/pflag"

var k8sOptions K8SOptions

type K8SOptions struct {
	KubeConfig string
}

func GetK8SFlagSet() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&k8sOptions.KubeConfig, "kube-config", "", "KubeConfig path")
	return flagSet
}

func GetK8SOptions() K8SOptions {
	return k8sOptions
}
