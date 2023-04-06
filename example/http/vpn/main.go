package main

import (
	"os"

	example "github.com/kube-peering/example"
	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/tunnel"
)

var (
	mode     tunnel.TunnelMode
	protocol = "http"
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
	server := tunnel.NewTunnelServer(mode, protocol, example.TunnelPort, example.TunnelServerCert, example.TunnelServerKey, example.TunnelServerName)

	if mode == tunnel.Forward {
		fowarder := pkg.NewForwarder(protocol, example.AppAddr)
		server.SetOnHTTPTunnel(fowarder.ForwardHTTP)
	}

	if mode == tunnel.Reverse {
		interceptor := pkg.NewInterceptor(protocol, example.VPNPort)
		interceptor.OnHTTPConnected = server.TunnelHTTP
		go interceptor.Start()
	}

	go server.Start()
}

func client() {
	client := tunnel.NewTunnelClient(mode, protocol, example.TunnelAddr, example.TunnelCaCert, example.TunnelServerName)

	if mode == tunnel.Forward {
		interceptor := pkg.NewInterceptor(protocol, example.VPNPort)
		interceptor.OnHTTPConnected = client.TunnelHTTP
		go interceptor.Start()
	}

	if mode == tunnel.Reverse {
		fowarder := pkg.NewForwarder(protocol, example.AppAddr)
		client.SetOnHTTPTunnel(fowarder.ForwardHTTP)
	}

	go client.Start()
}
