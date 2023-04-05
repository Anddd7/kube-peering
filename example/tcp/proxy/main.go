package main

import (
	"github.com/kube-peering/internal/pkg"
)

func main() {
	interceptor := pkg.NewInterceptor("tcp", 10021)
	forwarder := pkg.NewFowarder("tcp", ":8080")
	interceptor.OnTCPConnected = forwarder.ForwardTCP

	interceptor.Start()
}
