package cmd

import (
	"github.com/kube-peering/internal/kpctl"
	"github.com/kube-peering/internal/kpeering"
	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/spf13/cobra"
)

var instance *kpeering.Kpeering

var startCmd = &cobra.Command{
	Use: "start",
	PreRun: func(cmd *cobra.Command, args []string) {
		_logger := logger.CreateLocalLogger().With(
			"cmd", "kpeering start",
		)
		cfg, err := kpctl.ReadConfig(config.ConfigFile)
		if err != nil {
			_logger.Errorf("read config file failed: %s", err)
			return
		}

		instance = &kpeering.Kpeering{
			Logger:    _logger,
			VPNConfig: cfg.VPNConfig,
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		instance.Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
