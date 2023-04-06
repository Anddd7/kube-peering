package main

import (
	"os"

	example "github.com/kube-peering/example"
	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/tunnel"
)

func main() {
	server()
	client()

	select {}
}

func server() {
	_, port := tunnelPorts()

	var t pkg.Tunnel
	if port == 0 {
		// vpn tunnel
		t = tunnel.NewTunnelServer(pkg.Forward, "http", 10086, example.TunnelServerCert, example.TunnelServerKey, example.TunnelServerName)
		fowarder := pkg.NewFowarder("http", "localhost:8080")
		t.SetOnHTTPTunnel(fowarder.ForwardHTTP)
	} else {
		// reverse tunnel
		t = tunnel.NewTunnelServer(pkg.Reverse, "http", 10086, example.TunnelServerCert, example.TunnelServerKey, example.TunnelServerName)
		interceptor := pkg.NewInterceptor("http", port)
		interceptor.OnHTTPConnected = t.TunnelHTTP
		go interceptor.Start()
	}
	go t.Start()
}

func client() {
	port, _ := tunnelPorts()

	var t pkg.Tunnel
	if port == 0 {
		// reverse tunnel
		t = tunnel.NewTunnelClient(pkg.Reverse, "http", "localhost:10086", example.TunnelCaCert, example.TunnelServerName)
		fowarder := pkg.NewFowarder("http", "localhost:8080")
		t.SetOnHTTPTunnel(fowarder.ForwardHTTP)
	} else {
		// vpn tunnel
		t = tunnel.NewTunnelClient(pkg.Forward, "http", "localhost:10086", example.TunnelCaCert, example.TunnelServerName)
		interceptor := pkg.NewInterceptor("http", port)
		interceptor.OnHTTPConnected = t.TunnelHTTP
		go interceptor.Start()
	}
	go t.Start()
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
