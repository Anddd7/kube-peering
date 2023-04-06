package main

import (
	"os"

	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/tunnel"

	example "github.com/kube-peering/example"
)

var (
	mode     tunnel.TunnelMode
	protocol = "tcp"
)

func init() {
	if len(os.Args) > 1 {
		if os.Args[1] == "reverse" {
			mode = tunnel.Reverse
		}
	}
	mode = tunnel.Forward
}

func main() {
	server()
	client()

	select {}
}

func server() {
	server := tunnel.NewTunnelServer(mode, protocol, example.TunnelPort, example.TunnelServerCert, example.TunnelServerKey, example.TunnelServerName)

	if mode == tunnel.Forward {
		fowarder := pkg.NewFowarder(protocol, example.AppAddr)
		server.SetOnTCPTunnel(fowarder.ForwardTCP)
	}

	if mode == tunnel.Reverse {
		interceptor := pkg.NewInterceptor(protocol, example.VPNPort)
		interceptor.OnTCPConnected = server.TunnelTCP
		go interceptor.Start()
	}

	go server.Start()
}

func client() {
	client := tunnel.NewTunnelClient(mode, protocol, example.TunnelAddr, example.TunnelCaCert, example.TunnelServerName)

	if mode == tunnel.Forward {
		interceptor := pkg.NewInterceptor(protocol, example.VPNPort)
		interceptor.OnTCPConnected = client.TunnelTCP
		go interceptor.Start()
	}

	if mode == tunnel.Reverse {
		fowarder := pkg.NewFowarder(protocol, example.AppAddr)
		client.SetOnTCPTunnel(fowarder.ForwardTCP)
	}

	go client.Start()
}
