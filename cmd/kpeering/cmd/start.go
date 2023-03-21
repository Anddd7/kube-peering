package cmd

import (
	"net"

	"github.com/kube-peering/internal/io"
	"github.com/kube-peering/internal/logger"
	"github.com/kube-peering/internal/model"
	"github.com/spf13/cobra"
)

type Kpeering struct {
	frontdoor model.Frontdoor
	backdoor  model.Backdoor
}

var kpeering *Kpeering

var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitLogger()
		// TODO create via command args
		kpeering = &Kpeering{
			frontdoor: model.DefaultFrontdoor,
			backdoor:  model.DefaultBackdoor,
		}
		start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func start() {
	frontdoorListener, err := net.Listen("tcp", kpeering.frontdoor.Address())
	if err != nil {
		logger.Z.Errorf("failed to start frontdoor listener: %v", err)
		return
	}
	defer frontdoorListener.Close()

	backdoorListener, err := net.Listen("tcp", kpeering.backdoor.Address())
	if err != nil {
		logger.Z.Errorf("failed to start backdoor listener: %v", err)
		return
	}
	defer backdoorListener.Close()

	// TODO use mutex and wait for both front and back conn ready
	for {
		frontdoorConn, err := frontdoorListener.Accept()
		if err != nil {
			logger.Z.Error(err)
			continue
		}
		logger.Z.Infoln("frontdoor connection is comming")

		backdoorConn, err := backdoorListener.Accept()
		if err != nil {
			logger.Z.Error(err)
			continue
		}
		logger.Z.Infoln("backdoor is opened")

		go io.BiFoward("Frontdoor", frontdoorConn, "Backdoor", backdoorConn)
	}
}
