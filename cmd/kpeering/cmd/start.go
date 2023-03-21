package cmd

import (
	"github.com/kube-peering/internal/kpeering"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/model"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitLogger()
		// TODO create via command args
		kpeering := &kpeering.Kpeering{
			Frontdoor: model.DefaultFrontdoor,
			Backdoor:  model.DefaultBackdoor,
		}
		kpeering.Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
