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
				Backdoor:    model.DefaultBackdoor,
				Application: model.CreateApplication("localhost", 8080),
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			instance.Connect()
		},
	}
)

func init() {
	rootCmd.AddCommand(connectCmd)
}
