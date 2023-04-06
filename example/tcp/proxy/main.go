package main

import (
	example "github.com/kube-peering/example"
	"github.com/kube-peering/internal/pkg"
)

func main() {
	interceptor := pkg.NewInterceptor("tcp", example.ProxyPort)
	forwarder := pkg.NewFowarder("tcp", example.AppAddr)
	interceptor.OnTCPConnected = forwarder.ForwardTCP

	interceptor.Start()
}
