package main

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/spf13/cobra"
	"github.com/stevenayers/netkat"
	"os"
	"os/user"
)

var (
	config  string
	context string
)

var rootCmd = &cobra.Command{
	Use:   "netkat [TARGET URL]",
	Short: "Netkat is a CLI for troubleshooting kubernetes networking issues",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a url target")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		netkat.InitLogger(log.NewSyncWriter(os.Stdout), "error")
		var ch netkat.Checker
		if config == "" {
			usr, _ := user.Current()
			config = fmt.Sprintf("%v/.kube/config", usr.HomeDir)
		}
		ch.Client = netkat.InitClient(context, config)
		ch.KubernetesComponents = ch.Client.GetComponents()

		err := ch.ParseTarget(args[0])
		if err != nil {
			_ = level.Error(netkat.Logger).Log("msg", err)
		}
		ch.RunChecks()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&config, "config", "", "Kubernetes config file (default is $HOME/.kube/config)")
	rootCmd.PersistentFlags().StringVar(&context, "context", "default", "Kubernetes cluster context name")
}

func main() {
	Execute()
}
