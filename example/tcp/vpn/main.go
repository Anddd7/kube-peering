package main

import (
	"os"

	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/connectors"

	example "github.com/kube-peering/example"
)

var (
	mode pkg.TunnelMode
	cfg  = connectors.VPNConfig{
		Protocol:   pkg.TCP,
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
