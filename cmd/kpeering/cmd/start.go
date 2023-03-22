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
		instance = &kpeering.Kpeering{
			Frontdoor: model.CreateFrontdoor(flags.protocol, flags.host, flags.port),
			Backdoor:  model.DefaultBackdoor,
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		instance.Start()
	},
}

var flags = struct {
	protocol string
	host     string
	port     int
}{}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&flags.protocol, "protocol", "tcp", "the target protocol")
	startCmd.Flags().StringVar(&flags.host, "host", "localhost", "the target host")
	startCmd.Flags().IntVarP(&flags.port, "port", "p", config.DefautlFrontdoorPort, "the target port")
}
