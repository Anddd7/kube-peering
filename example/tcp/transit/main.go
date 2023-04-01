package main

import (
	"net"

	"github.com/kube-peering/internal/pkg/transit"
)

func main() {
	interceptor := transit.NewInterceptor("tcp", 10021)
	forwarder := transit.NewFowarder("tcp", ":8080")
	interceptor.OnTCPConnected = func(conn *net.TCPConn) {
		forwarder.ForwardTCP(conn)
	}

	interceptor.Start()
}
