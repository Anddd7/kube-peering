package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/kube-peering/internal/config"
	"github.com/kube-peering/internal/logger"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use: "connect",
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitLogger()
		connect()
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}

func connect() {
	kpeeringConn, err := net.Dial("tcp", config.KpeeringBackdoorAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to server: %v\n", err)
		os.Exit(1)
	}

	// TODO should get app address from config
	appConn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to app: %v\n", err)
		os.Exit(1)
	}

	forward(kpeeringConn, appConn)
}

func forward(from net.Conn, to net.Conn) {
	defer from.Close()
	defer to.Close()

	ch1 := make(chan []byte)
	ch2 := make(chan []byte)
	// read data and put into channel
	go func() {
		for {
			buf := make([]byte, 1024)
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
		for {
			buf := make([]byte, 1024)
			n, err := to.Read(buf)
			if err != nil {
				logger.Z.Errorf("Got and error: %v", err)
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
