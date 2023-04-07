package cmd

import (
	"strconv"
	"strings"

	"github.com/kube-peering/internal/kpctl"
	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/connectors"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	configureCmd = &cobra.Command{
		Use: "configure",
		Run: func(cmd *cobra.Command, args []string) {
			_logger := logger.CreateLocalLogger().With(
				"cmd", "kpctl configure",
			)
			if kpctl.IsConfigExists(config.ConfigFile) {
				_logger.Infof("config file %s exists.", config.ConfigFile)
				return
			}

			tunnelAddr := strings.Split(flags.tunnelAddr, ":")
			port, _ := strconv.Atoi(tunnelAddr[1])

			instance := kpctl.Config{
				VPNConfig: connectors.VPNConfig{
					Protocol: pkg.Protocol(flags.protocol),
					Tunnel: pkg.TunnelConfig{
						Host:           tunnelAddr[0],
						Port:           port,
						ServerCertPath: flags.tunnelServerCertPath,
						ServerKeyPath:  flags.tunnelServerKeyPath,
						CaCertPath:     flags.tunnelCACertPath,
						ServerName:     flags.tunnelServerName,
					},
					LocalPort:  flags.localPort,
					RemotePort: flags.remotePort,
				},
			}

			err := instance.WriteFile(config.ConfigFile)
			if err != nil {
				_logger.Errorf("write config file failed: %s", err)
				return
			}
		},
	}
)

var flags = struct {
	protocol             string
	tunnelAddr           string
	tunnelServerCertPath string
	tunnelServerKeyPath  string
	tunnelCACertPath     string
	tunnelServerName     string
	localPort            int
	remotePort           int
}{}

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.Flags().StringVar(&flags.protocol, "protocol", "tcp", "the protocol of peer connection")
	configureCmd.Flags().StringVar(&flags.tunnelAddr, "tunnel", "example.com:10022", "the tunnel server address")
	configureCmd.Flags().StringVar(&flags.tunnelServerCertPath, "tunnel-server-cert", "./kpctl/server.crt", "the server cert of tunnel server")
	configureCmd.Flags().StringVar(&flags.tunnelServerKeyPath, "tunnel-server-key", "./kpctl/server.key", "the server key of tunnel server")
	configureCmd.Flags().StringVar(&flags.tunnelCACertPath, "tunnel-ca-cert", "./kpctl/ca.crt", "the ca cert of tunnel server")
	configureCmd.Flags().StringVar(&flags.tunnelServerName, "tunnel-server-name", "example.com", "the server name of tunnel server")
	configureCmd.Flags().IntVar(&flags.localPort, "local-port", 8080, "the port of application running on local")
	configureCmd.Flags().IntVar(&flags.remotePort, "remote-port", 8080, "the port of application running remotely")
}
