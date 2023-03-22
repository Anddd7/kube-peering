/*
Copyright © 2023 Anddd7 <liaoad_space@sina.com>
*/
package cmd

import (
	"os"

	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "kpeering",
	Short: "manage the proxy server in container/cluster",
	Long:  `...`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logger.Z.Error(`failed to execute command`, zap.Error(err))
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&config.DebugMode, "debug", config.DebugMode, "Enable debug logs")
	rootCmd.Flags().StringVar(&config.LogEncoder, "log-encoder", config.LogEncoder, "Log format, json or plain")
}
