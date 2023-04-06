package main

import (
	"os"

	example "github.com/kube-peering/example"
	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/connectors"
	"github.com/kube-peering/internal/pkg/tunnel"
)

var (
	mode tunnel.TunnelMode
	cfg  = connectors.VPNConfig{
		Protocol:   pkg.HTTP,
		LocalPort:  example.VPNPort,
		RemoteAddr: example.AppAddr,
		Tunnel: connectors.TunnelConfig{
			Port:           example.TunnelPort,
			Host:           example.TunnelAddr,
			ServerCertPath: example.TunnelServerCert,
			ServerKeyPath:  example.TunnelServerKey,
			CaCertPath:     example.TunnelCaCert,
			ServerName:     example.TunnelServerName,
		},
	}
)

func init() {
	if len(os.Args) > 1 {
		if os.Args[1] == "reverse" {
			mode = tunnel.Reverse
		}
	}
	mode = tunnel.Forward
}

// client will connect to 10022
// normal : client -> tunnel client --------> tunnel server -> server
// reverse: client -> tunnel server --------> tunnel client -> server
func main() {
	server()
	client()

	select {}
}

func server() {
	if mode == tunnel.Forward {
		connectors.NewVPNServer(cfg).Start()
	}

	if mode == tunnel.Reverse {
		connectors.NewPortFowardServer(cfg).Start()
	}
}

func client() {
	if mode == tunnel.Forward {
		connectors.NewPortFowardServer(cfg).Start()
	}

	if mode == tunnel.Reverse {
		connectors.NewPortForwardClient(cfg).Start()
	}
}
