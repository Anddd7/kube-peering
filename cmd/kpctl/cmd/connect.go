package cmd

import (
	"net"
	"os"

	"github.com/kube-peering/internal/io"
	"github.com/kube-peering/internal/logger"
	"github.com/kube-peering/internal/model"
	"github.com/spf13/cobra"
)

type Kpctl struct {
	backdoor    model.Backdoor
	application model.Application
}

var kpctl *Kpctl

var connectCmd = &cobra.Command{
	Use: "connect",
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitLogger()
		// TODO create via command args
		kpctl = &Kpctl{
			backdoor:    model.DefaultBackdoor,
			application: model.CreateApplication("localhost", 8080),
		}
		connect()
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}

func connect() {
	backdoorConn, err := net.Dial("tcp", kpctl.backdoor.Address())
	if err != nil {
		logger.Z.Errorf("failed to connect backdoor of peering server: %v", err)
		os.Exit(1)
	}

	applicationConn, err := net.Dial("tcp", kpctl.application.Address())
	if err != nil {
		logger.Z.Errorf("failed to connect target application: %v", err)
		os.Exit(1)
	}

	io.BiFoward("Backdoor", backdoorConn, "Application", applicationConn)
}
