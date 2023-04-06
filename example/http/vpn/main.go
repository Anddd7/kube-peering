package main

import (
	"os"

	example "github.com/kube-peering/example"
	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/connectors"
)

var (
	mode pkg.TunnelMode
	cfg  = connectors.VPNConfig{
		Protocol:   pkg.HTTP,
		LocalPort:  example.VPNPort,
		RemoteAddr: example.AppAddr,
		Tunnel: pkg.TunnelConfig{
			Port:           example.TunnelPort,
			Host:           example.TunnelHost,
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
			mode = pkg.Reverse
		}
	}
	mode = pkg.Forward
}

// client will connect to 10022
// normal : client -> tunnel client --------> tunnel server -> server
// reverse: client -> tunnel server --------> tunnel client -> server
func main() {
	go server()
	go client()

	select {}
}

func server() {
	if mode == pkg.Forward {
		connectors.NewVPNServer(cfg).Start()
	}

	if mode == pkg.Reverse {
		connectors.NewReverseVPNServer(cfg).Start()
	}
}

func client() {
	if mode == pkg.Forward {
		connectors.NewVPNClient(cfg).Start()
	}

	if mode == pkg.Reverse {
		connectors.NewReverseVPNClient(cfg).Start()
	}
}
