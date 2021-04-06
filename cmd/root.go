package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/qumine/ingress-controller/internal/api"
	"github.com/qumine/ingress-controller/internal/ingress"
	"github.com/qumine/ingress-controller/internal/k8s"
	"github.com/qumine/ingress-controller/pkg/build"
	"github.com/qumine/ingress-controller/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Execute() {
	err := NewRootCmd().Execute()
	if err != nil {
		logrus.Fatal(err)
	}
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "ingress-controller",
		Short:   "A Kubernetes ingress controller for minecraft servers",
		Long:    "A Kubernetes ingress controller for minecraft servers",
		Version: build.Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			viper.AutomaticEnv()
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				if !f.Changed && viper.IsSet(f.Name) {
					val := viper.Get(f.Name)
					cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
				}
			})

			cliOptions := config.GetCliOptions()
			logrus.SetLevel(cliOptions.LogLevel)
		},
		Run: func(cmd *cobra.Command, args []string) {
			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
			ctx, cancel := context.WithCancel(context.Background())
			wg := &sync.WaitGroup{}

			k8s := k8s.NewK8S(config.GetK8SOptions())
			ing := ingress.NewIngress(config.GetIngressOptions())
			api := api.NewAPI(config.GetAPIOptions(), k8s, ing)

			go k8s.Start(ctx, wg)
			go ing.Start(ctx, wg)
			go api.Start(ctx, wg)

			<-interrupt
			logrus.Info("Interrupted, stopping")

			cancel()
			wg.Wait()
		},
	}
	rootCmd.PersistentFlags().AddFlagSet(config.GetCliFlagSet())
	rootCmd.PersistentFlags().AddFlagSet(config.GetK8SFlagSet())
	rootCmd.PersistentFlags().AddFlagSet(config.GetAPIFlagSet())
	rootCmd.PersistentFlags().AddFlagSet(config.GetIngressFlagSet())
	return rootCmd
}
