package main

import (
	"os"

	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/tunnel"

	example "github.com/kube-peering/example"
)

func main() {
	server()
	client()

	select {}
}

func server() {
	_, port := tunnelPorts()
	tunnel := tunnel.NewTunnelServer(pkg.Forward, "tcp", 10086, example.TunnelServerCert, example.TunnelServerKey, example.TunnelServerName)

	if port == 0 {
		fowarder := pkg.NewFowarder("tcp", ":8080")
		tunnel.SetOnTCPTunnel(fowarder.ForwardTCP)
	} else {
		interceptor := pkg.NewInterceptor("tcp", port)
		interceptor.OnTCPConnected = tunnel.TunnelTCP
		go interceptor.Start()
	}

	go tunnel.Start()
}

func client() {
	port, _ := tunnelPorts()
	tunnel := tunnel.NewTunnelClient(pkg.Forward, "tcp", "localhost:10086", example.TunnelCaCert, example.TunnelServerName)

	if port == 0 {
		fowarder := pkg.NewFowarder("tcp", ":8080")
		tunnel.SetOnTCPTunnel(fowarder.ForwardTCP)
	} else {
		interceptor := pkg.NewInterceptor("tcp", port)
		interceptor.OnTCPConnected = tunnel.TunnelTCP
		go interceptor.Start()
	}

	go tunnel.Start()
}

// client will connect to 10022
// normal : client -> tunnel client --------> tunnel server -> server
// reverse: client -> tunnel server --------> tunnel client -> server
func tunnelPorts() (int, int) {
	if len(os.Args) > 1 {
		if os.Args[1] == "reverse" {
			return 0, 10022
		}
	}
	return 10022, 0
}
