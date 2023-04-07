package cmd

import (
	"net"

	"github.com/kube-peering/internal/kpctl"
	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/connectors"
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
			tunnelAddr, _ := net.ResolveTCPAddr("tcp", flags.tunnelAddr)

			instance = &kpctl.Kpctl{
				Logger: _logger,
				VPNConfig: connectors.VPNConfig{
					Protocol:  pkg.Protocol(flags.protocol),
					LocalPort: flags.localPort,
					Tunnel: pkg.TunnelConfig{
						Host:       tunnelAddr.IP.String(),
						Port:       tunnelAddr.Port,
						CaCertPath: flags.tunnelCACert,
						ServerName: flags.tunnelServerName,
					},
				},
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			instance.Connect()
		},
	}
)

// TODO use toml config file to replace long flags
var flags = struct {
	protocol         string
	tunnelAddr       string
	tunnelCACert     string
	tunnelServerName string
	localPort        int
	remotePort       int
}{}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringVar(&flags.protocol, "protocol", "tcp", "the protocol of peer connection")
	connectCmd.Flags().StringVar(&flags.tunnelAddr, "tunnel", "localhost:10022", "the tunnel server address")
	connectCmd.Flags().StringVar(&flags.tunnelCACert, "tunnel-ca-cert", "../../bin/ca.crt", "the ca cert of tunnel server")
	connectCmd.Flags().StringVar(&flags.tunnelServerName, "tunnel-server-name", "localhost", "the tunnel server nam in cert")
	connectCmd.Flags().IntVar(&flags.localPort, "local-port", 8080, "the port of application running on local")
	connectCmd.Flags().IntVar(&flags.remotePort, "remote-port", 8080, "the port of application running remotely")
}
