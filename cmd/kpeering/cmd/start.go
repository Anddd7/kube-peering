package cmd

import (
	"net"

	"github.com/kube-peering/internal/config"
	"github.com/kube-peering/internal/logger"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitLogger()
		start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func forward(from net.Conn, to net.Conn) {
	defer from.Close()
	defer to.Close()

	ch1 := make(chan []byte)
	ch2 := make(chan []byte)
	// read data and put into channel
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := from.Read(buf)
			if err != nil {
				logger.Z.Error(err)
				break
			}
			logger.Z.Infoln("Recive msg from client side")
			ch1 <- buf[:n]
		}
	}()
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := to.Read(buf)
			if err != nil {
				logger.Z.Error(err)
				break
			}
			logger.Z.Infoln("Recive msg from server side")
			ch2 <- buf[:n]
		}
	}()

	for {
		select {
		case buf := <-ch1:
			_, err := to.Write(buf)
			logger.Z.Infoln("Write msg to backdoor side")
			if err != nil {
				logger.Z.Error(err)
				break
			}
		case buf := <-ch2:
			_, err := from.Write(buf)
			logger.Z.Infoln("Write msg to foward side")
			if err != nil {
				logger.Z.Error(err)
				break
			}
		}
	}
}

func start() {
	forwardListener, err := net.Listen("tcp", config.KpeeringForwardAddr)
	if err != nil {
		logger.Z.Error(err)
		return
	}
	defer forwardListener.Close()

	backdoorListener, err := net.Listen("tcp", config.KpeeringBackdoorAddr)
	if err != nil {
		logger.Z.Error(err)
		return
	}
	defer backdoorListener.Close()

	for {
		backdoorConn, err := backdoorListener.Accept()
		if err != nil {
			logger.Z.Error(err)
			continue
		}
		logger.Z.Infoln("backdoor is open")

		forwardConn, err := forwardListener.Accept()
		if err != nil {
			logger.Z.Error(err)
			continue
		}
		logger.Z.Infoln("accept forward request")

		go forward(forwardConn, backdoorConn)
	}
}
