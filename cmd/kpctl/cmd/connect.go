package cmd

import (
	"github.com/kube-peering/internal/kpctl"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/model"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use: "connect",
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitLogger()
		// TODO create via command args
		kpctl := &kpctl.Kpctl{
			Backdoor:    model.DefaultBackdoor,
			Application: model.CreateApplication("localhost", 8080),
		}
		kpctl.Connect()
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
