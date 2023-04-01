package cmd

import (
	"net"
	"strconv"

	"github.com/kube-peering/internal/kpctl"
	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/model"
	"github.com/spf13/cobra"
)

var (
	instance   *kpctl.Kpctl
	connectCmd = &cobra.Command{
		Use: "connect",
		PreRun: func(cmd *cobra.Command, args []string) {
			logger.InitLogger(config.DebugMode, config.LogEncoder)

			instance = &kpctl.Kpctl{
				Tunnel:    createTunnelClient(),
				Forwarder: model.CreateForwarder("localhost", 8080),
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			instance.Connect()
		},
	}
)

func createTunnelClient() model.Tunnel {
	host, _port, err := net.SplitHostPort(flags.tunnelAddr)
	if err != nil {
		panic(err)
	}
	port, err := strconv.Atoi(_port)
	if err != nil {
		panic(err)
	}

	if flags.protocol == "tcp" {
		return model.Tunnel{
			Endpoint: model.Endpoint{
				Name:       "tcp-tunnel",
				Protocol:   "tcp",
				Host:       host,
				ListenPort: port,
			},
		}
	}

	return model.Tunnel{
		Endpoint: model.Endpoint{
			Name:       flags.protocol + "-tunnel",
			Protocol:   flags.protocol,
			Host:       host,
			ListenPort: port,
		},
	}
}

// TODO use toml config file to replace long flags
var flags = struct {
	protocol         string
	tunnelAddr       string
	tunnelCACert     string
	tunnelServerName string
	applicationAddr  string
}{}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringVar(&flags.protocol, "protocol", "tcp", "the protocol of peer connection")
	connectCmd.Flags().StringVar(&flags.tunnelAddr, "tunnel", "localhost:10022", "the tunnel server address")
	connectCmd.Flags().StringVar(&flags.tunnelCACert, "tunnel-ca-cert", "../../bin/ca.crt", "the ca cert of tunnel server")
	connectCmd.Flags().StringVar(&flags.tunnelServerName, "tunnel-server-name", "localhost", "the tunnel server nam in cert")
	connectCmd.Flags().StringVar(&flags.applicationAddr, "application", "localhost:8080", "the proxied application address")
}
