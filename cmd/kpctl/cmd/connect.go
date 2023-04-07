package cmd

import (
	"github.com/kube-peering/internal/kpctl"
	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	instance   *kpctl.Kpctl
	connectCmd = &cobra.Command{
		Use: "connect",
		PreRun: func(cmd *cobra.Command, args []string) {
			_logger := logger.CreateLocalLogger().With(
				"cmd", "kpctl connect",
			)

			cfg, err := kpctl.ReadConfig(config.ConfigFile)
			if err != nil {
				_logger.Errorf("read config file failed: %s", err)
				return
			}
			instance = &kpctl.Kpctl{
				Logger:    _logger,
				VPNConfig: cfg.VPNConfig,
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
