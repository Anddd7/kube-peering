package cmd

import (
	"github.com/kube-peering/internal/kpeering"
	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/connectors"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/spf13/cobra"
)

var instance *kpeering.Kpeering

var startCmd = &cobra.Command{
	Use: "start",
	PreRun: func(cmd *cobra.Command, args []string) {
		_logger := logger.CreateLocalLogger().With(
			"cmd", "kpeering start",
		)

		instance = &kpeering.Kpeering{
			Logger: _logger,
			VPNConfig: connectors.VPNConfig{
				Protocol:   pkg.Protocol(flags.protocol),
				RemotePort: flags.port,
				Tunnel: pkg.TunnelConfig{
					Host:           "localhost",
					Port:           flags.tunnelPort,
					ServerCertPath: flags.tunnelServerCert,
					ServerKeyPath:  flags.tunnelServerKey,
					ServerName:     flags.tunnelServerName,
				},
			},
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		instance.Start()
	},
}

var flags = struct {
	protocol         string
	tunnelPort       int
	tunnelServerCert string
	tunnelServerKey  string
	tunnelServerName string
	port             int
}{}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&flags.protocol, "protocol", "tcp", "the protocol of peer connection")
	startCmd.Flags().IntVar(&flags.tunnelPort, "tunnel-port", 10022, "the tunnel server port")
	startCmd.Flags().StringVar(&flags.tunnelServerCert, "tunnel-server-cert", "../../bin/server.crt", "the server cert of tunnel server")
	startCmd.Flags().StringVar(&flags.tunnelServerKey, "tunnel-server-key", "../../bin/server.key", "the server key of tunnel server")
	startCmd.Flags().StringVar(&flags.tunnelServerName, "tunnel-server-name", "localhost", "the tunnel server nam in cert")
	startCmd.Flags().IntVarP(&flags.port, "port", "p", 8080, "the port of application running remotely")
}
