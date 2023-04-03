package main

import (
	"crypto/tls"

	"github.com/kube-peering/internal/pkg/transit"

	example "github.com/kube-peering/example"
)

func main() {
	server()
	client()

	select {}
}

func server() {
	// interceptor := transit.NewInterceptor("tcp", 10021)
	tunnel := transit.NewTunnelServer("tcp", 10086, example.TunnelServerCert, example.TunnelServerKey, example.TunnelServerName)
	fowarder := transit.NewFowarder("tcp", ":8080")

	tunnel.OnTlsConnected = func(conn *tls.Conn) {
		fowarder.ForwardTls(conn)
	}

	go tunnel.Start()
}

func client() {
	interceptor := transit.NewInterceptor("tcp", 10022)
	tunnel := transit.NewTunnelClient("localhost:10086", example.TunnelCaCert, example.TunnelServerName)

	interceptor.OnTCPConnected = tunnel.ForwardTls

	go interceptor.Start()
	go tunnel.Start()
}
