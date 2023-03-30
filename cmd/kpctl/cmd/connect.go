package cmd

import (
	"github.com/kube-peering/internal/kpctl"
	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/model"
	"github.com/spf13/cobra"
)

var (
	instance   *kpctl.Kpctl
	connectCmd = &cobra.Command{
		Use: "connect",
		PreRun: func(cmd *cobra.Command, args []string) {
			logger.InitLogger(config.DebugMode, config.LogEncoder)
			instance = &kpctl.Kpctl{
				Tunnel:    model.DefaultTunnel,
				Forwarder: model.CreateForwarder("localhost", 8080),
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			instance.Connect()
		},
	}
)

var flags = struct {
	tcp  bool
	http bool
	port int
}{}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().BoolVar(&flags.tcp, "tcp", true, "build a tcp tunnel")
	connectCmd.Flags().BoolVar(&flags.http, "http", false, "build a http tunnel")
	connectCmd.Flags().IntVarP(&flags.port, "port", "p", 0, "the listening port")
	err := connectCmd.MarkFlagRequired("port")
	if err != nil {
		panic(err)
	}
}
