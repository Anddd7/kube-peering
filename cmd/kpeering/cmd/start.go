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
			Frontdoor: model.DefaultFrontdoor,
			Backdoor:  model.DefaultBackdoor,
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		instance.Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
