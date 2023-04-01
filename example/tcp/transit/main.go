package main

import (
	"fmt"
	"net"

	"github.com/kube-peering/internal/pkg/transit"
)

func main() {
	interceptor := transit.NewInterceptor("tcp", 10021)
	proxy := transit.NewProxy("tcp", ":8080")
	interceptor.OnTCPConnected = func(conn *net.TCPConn) {
		proxy.ProxyTCP(conn)
	}

	interceptor.Start()
}
