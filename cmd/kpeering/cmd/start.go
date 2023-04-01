package cmd

import (
	"github.com/kube-peering/internal/kpeering"
	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/model"
	"github.com/spf13/cobra"
)

var instance *kpeering.Kpeering

var startCmd = &cobra.Command{
	Use: "start",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.InitLogger(config.DebugMode, config.LogEncoder)
		protocol := "tcp"
		if flags.http {
			protocol = "http"
		}
		instance = &kpeering.Kpeering{
			Interceptor: model.CreateInterceptor(protocol, flags.port),
			Tunnel: model.CreateTunnelServer(
				"localhost", 10022,
				"../../bin/server.key",
				"../../bin/server.crt",
				"localhost",
			),
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		instance.Start()
	},
}

var flags = struct {
	tcp  bool
	http bool
	port int
}{}

//lintignore:errorcheck
func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolVar(&flags.tcp, "tcp", true, "build a tcp tunnel")
	startCmd.Flags().BoolVar(&flags.http, "http", false, "build a http tunnel")
	startCmd.Flags().IntVarP(&flags.port, "port", "p", 0, "the listening port")
	err := startCmd.MarkFlagRequired("port")
	if err != nil {
		panic(err)
	}
}
