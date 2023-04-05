package main

import (
	"net"

	"github.com/kube-peering/internal/pkg"
)

func main() {
	interceptor := pkg.NewInterceptor("tcp", 10021)
	forwarder := pkg.NewFowarder("tcp", ":8080")
	interceptor.OnTCPConnected = func(conn *net.TCPConn) {
		forwarder.ForwardTCP(conn)
	}

	interceptor.Start()
}
